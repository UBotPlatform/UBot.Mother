package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-http-utils/etag"
	"github.com/gorilla/mux"
)

func ServeWebUI() {
	if Config.WebUI.Address == "" {
		return
	}
	addr := Config.WebUI.Address
	fmt.Println("WebUI listening at:", addr)
	r := mux.NewRouter()

	r.Handle("/api/mother/services", etag.Handler(http.HandlerFunc(getServices), false)).Methods(http.MethodGet)
	r.HandleFunc("/api/mother/services", putServices).Methods(http.MethodPut)
	r.Handle("/api/mother/accounts", etag.Handler(http.HandlerFunc(getAccounts), false)).Methods(http.MethodGet)
	r.HandleFunc("/api/mother/accounts", postAccounts).Methods(http.MethodPost)
	r.HandleFunc("/api/mother/account_providers", getAccountProviders).Methods(http.MethodGet)
	r.Handle("/api/mother/apps", etag.Handler(http.HandlerFunc(getApps), false)).Methods(http.MethodGet)
	r.HandleFunc("/api/mother/apps/reload", reloadApps).Methods(http.MethodPost)
	r.PathPrefix("/").Handler(WebUIStaticHandler)
	r.Use(BasicAuthMiddleware(Config.WebUI.User, Config.WebUI.Password))
	err := http.ListenAndServe(addr, CORSMiddleware(r))
	if err != nil {
		log.Println("WebUI cannot listen or serve:", err)
	}
}

func getServices(writer http.ResponseWriter, request *http.Request) {
	var r []*ServiceInfoFromMother
	for pair := Services.Oldest(); pair != nil; pair = pair.Next() {
		s := pair.Value
		r = append(r, buildServiceInfoResp(s))
	}
	webuiResponse(request, writer, r)
}

func unknownClientID(s *ServiceLaunchInfo) string {
	return fmt.Sprintf("Unknown@%x", sha1.New().Sum([]byte(s.Service.ID()))[:4])
}

func buildServiceInfoResp(s *ServiceLaunchInfo) *ServiceInfoFromMother {
	if s == nil {
		return nil
	}
	return &ServiceInfoFromMother{
		ID:       s.Service.ID(),
		Status:   s.Service.Status(),
		LaunchAt: s.LaunchAt,
	}
}

func buildServiceInfoRespByID(id string) *ServiceInfoFromMother {
	s, ok := Services.Get(id)
	if ok {
		return buildServiceInfoResp(s)
	}
	return nil
}

func buildServiceInfoRespByMatedate(matedate map[string]string) *ServiceInfoFromMother {
	if matedate == nil {
		return nil
	}
	if motherId, isChild := matedate["x-mother-id"]; isChild {
		if motherId == MotherID {
			return buildServiceInfoRespByID(matedate["x-service-id"])
		}
	}
	return nil
}

func getAccountProviders(writer http.ResponseWriter, request *http.Request) {
	var r []AccountProviderInfoFromMother
	for id := range AccountProviders {
		r = append(r, AccountProviderInfoFromMother{
			ID: id,
		})
	}
	webuiResponse(request, writer, r)
}
func getAccounts(writer http.ResponseWriter, request *http.Request) {
	var r []AccountInfoFromMother
	var list []AccountInfoFromRouter
	knowns := make(map[string]bool)
	if ManagerRPCConn == nil {
		http.Error(writer, "Router is not connected", http.StatusInternalServerError)
		return
	}
	err := ManagerRPCConn.CallExplicitly("get_account_list", []int{}, &list)
	if err != nil {
		webuiErrorResponse(request, writer, err)
		return
	}
	for _, x := range list {
		y := AccountInfoFromMother{
			ID:            x.ID,
			BindedService: buildServiceInfoRespByMatedate(x.ManagerMetadata),
		}
		r = append(r, y)
		if y.BindedService != nil {
			knowns[y.BindedService.ID] = true
		}
	}
	for pair := Services.Oldest(); pair != nil; pair = pair.Next() {
		s := pair.Value
		sid := s.Service.ID()
		if strings.HasPrefix(sid, "Account#") {
			if _, known := knowns[sid]; !known {
				r = append(r, AccountInfoFromMother{
					ID:            unknownClientID(s),
					BindedService: buildServiceInfoResp(s),
				})
			}
		}
	}
	webuiResponse(request, writer, r)
}

func getApps(writer http.ResponseWriter, request *http.Request) {
	var r []AppInfoFromMother
	var list []AppInfoFromRouter
	knowns := make(map[string]bool)
	if ManagerRPCConn == nil {
		http.Error(writer, "Router is not connected", http.StatusInternalServerError)
		return
	}
	err := ManagerRPCConn.CallExplicitly("get_app_list", []int{}, &list)
	if err != nil {
		webuiErrorResponse(request, writer, err)
		return
	}
	for _, x := range list {
		y := AppInfoFromMother{
			ID:            x.ID,
			BindedService: buildServiceInfoRespByMatedate(x.ManagerMetadata),
		}
		r = append(r, y)
		if y.BindedService != nil {
			knowns[y.BindedService.ID] = true
		}
	}
	for pair := Services.Oldest(); pair != nil; pair = pair.Next() {
		s := pair.Value
		sid := s.Service.ID()
		if strings.HasPrefix(sid, "App#") {
			if _, known := knowns[sid]; !known {
				r = append(r, AppInfoFromMother{
					ID:            unknownClientID(s),
					BindedService: buildServiceInfoResp(s),
				})
			}
		}
	}
	webuiResponse(request, writer, r)
}

func postAccounts(writer http.ResponseWriter, request *http.Request) {
	var r []SuccessResponse
	var data []AccountInfo
	err := webuiGetRequest(request, &data)
	if err != nil {
		webuiErrorResponse(request, writer, err)
		return
	}
	for _, item := range data {
		Config.Accounts = append(Config.Accounts, item)
		err = LoadAccount(item, true)
		if err == nil {
			r = append(r, SuccessResponse{true})
		} else {
			r = append(r, SuccessResponse{false})
		}
	}
	SaveConfig()
	webuiResponse(request, writer, r)
}

func putServices(writer http.ResponseWriter, request *http.Request) {
	var r []SuccessResponse
	var data []map[string]interface{}
	err := webuiGetRequest(request, &data)
	if err != nil {
		webuiErrorResponse(request, writer, err)
		return
	}
	for _, item := range data {
		success := false
		if id, ok := item["id"].(string); ok {
			if s, ok := Services.Get(id); ok {
				success = true
				if status, ok := item["status"].(float64); ok {
					switch int32(status) {
					case ServiceStarting:
						fallthrough
					case ServiceRunning:
						if s.Service.Status() == ServiceStopped {
							s.Start()
						}
					case ServiceStopped:
						if s.Service.Status() == ServiceRunning {
							s.Service.Stop()
						}
					}
				}
			}
		}
		r = append(r, SuccessResponse{success})
	}
	webuiResponse(request, writer, r)
}

func reloadApps(writer http.ResponseWriter, request *http.Request) {
	LoadApps(true)
	webuiResponse(request, writer, SuccessResponse{true})
}
func webuiGetRequest(request *http.Request, data interface{}) error {
	bin, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bin, data)
	return err
}
func webuiErrorResponse(request *http.Request, writer http.ResponseWriter, err error) {
	http.Error(writer, err.Error(), http.StatusInternalServerError)
}
func webuiResponse(request *http.Request, writer http.ResponseWriter, response interface{}) {
	jsonBinary, err := json.Marshal(response)
	if err != nil {
		webuiErrorResponse(request, writer, err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Length", strconv.Itoa(len(jsonBinary)))
	writer.Header().Set("Cache-control", "no-cache")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(jsonBinary)
}

func BasicAuthMiddleware(expectedUser string, expectedPassword string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			user, password, _ := request.BasicAuth()
			if (expectedUser != "" && expectedUser != user) ||
				(expectedPassword != "" && expectedPassword != password) {
				writer.Header().Add("WWW-Authenticate", "Basic realm=\"Secure Area\"")
				writer.Header().Add("Content-Type", "text/plain")
				writer.WriteHeader(http.StatusUnauthorized)
				_, _ = writer.Write([]byte("Unauthorized"))
				return
			}
			next.ServeHTTP(writer, request)
		})
	}
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token")
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
		writer.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, PUT, PATCH, DELETE")
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
