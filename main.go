package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/1354092549/wsrpc"
	"github.com/gorilla/websocket"
	orderedmap "github.com/wk8/go-ordered-map"
)

type ServiceLaunchInfo struct {
	Service    ClientService
	RetryIndex int32
	LaunchAt   time.Time
}

func (s *ServiceLaunchInfo) Start() {
	s.LaunchAt = time.Now()
	s.Service.Start(OnServiceExit, s)
}

var Services = orderedmap.New()
var AccountProviders map[string]ServiceInfo
var ServiceRootContext context.Context
var MotherID = NewToken()
var ManagerToken string
var ManagerRPCConn *wsrpc.WebsocketRPCConn
var AccountCounter int32 = -1
var OnServiceExit = make(chan interface{}, 5)

func FetchServiceInfo(dir string, prefix string) map[string]ServiceInfo {
	r := make(map[string]ServiceInfo)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return r
	}
	possibles := make(map[string]int)
	ubotSuffix := ".ubot"
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ubotSuffix) {
			possibles[strings.TrimSuffix(file.Name(), ubotSuffix)] = 0
		}
		if !file.IsDir() {
			nameWithoutExt := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))
			if strings.HasSuffix(nameWithoutExt, ubotSuffix) {
				possibles[strings.TrimSuffix(nameWithoutExt, ubotSuffix)] = 0
			}
		}
	}
clientLoop:
	for clientId := range possibles {
		name := prefix + clientId
		for _, ext := range PathExt {
			if FileExists(filepath.Join(dir, clientId+ubotSuffix+ext)) {
				r[name] = StandaloneProcessServiceInfo{
					Path: filepath.Join(dir, clientId+ubotSuffix+ext),
				}
				continue clientLoop
			}
		}
		log.Printf("failed to load client %s because the format is not supported", clientId)
	}
	return r
}

func main() {
	var cancelService context.CancelFunc
	ServiceRootContext, cancelService = context.WithCancel(context.Background())
	cancelSignal := make(chan os.Signal)
	signal.Notify(cancelSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cancelSignal
		cancelService()
		fmt.Println("Shutting down...")
	}()

	executableFile, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	defaultRootFolder := filepath.Dir(executableFile)
	flag.StringVar(&RootFolder, "root", defaultRootFolder, "")
	flag.Parse()
	err = FindUBotEnv()
	if err != nil {
		log.Fatal(err)
	}

	if runtime.GOOS == "windows" {
		_ = exec.Command("taskkill", "/F", "/T", "/IM", RouterImageName).Run()
	}

	LoadConfig()
	if Config.UBot.Address == "" {
		Config.UBot.Address = "localhost:5000"
	}
	SaveConfig()

	LoadApps(false)
	AccountProviders = FetchServiceInfo(AccountFolder, "")
	LoadAccounts(false)
	go ServeWebUI()

	DaemonThread(ServiceRootContext)
}
func LoadAccount(info AccountInfo, start bool) error {
	provider, providerExists := AccountProviders[info.Type]
	if !providerExists {
		return fmt.Errorf("cannot find provider for account type")
	}
	serviceInfo := provider
	name := fmt.Sprintf("Account#%d_%s", atomic.AddInt32(&AccountCounter, 1), info.Type)
	sli := &ServiceLaunchInfo{Service: NewService(serviceInfo, name, info.Args)}
	Services.Set(name, sli)
	if start {
		sli.Start()
	}
	return nil
}
func LoadAccounts(start bool) {
	for _, item := range Config.Accounts {
		err := LoadAccount(item, start)
		if err != nil {
			log.Printf("cannot load account (%s): %v", item.Type, err)
		}
	}
}
func LoadApps(start bool) {
	for name, serviceInfo := range FetchServiceInfo(AppFolder, "App#") {
		_, exists := Services.Get(name)
		if exists {
			continue
		}
		sli := &ServiceLaunchInfo{Service: NewService(serviceInfo, name, nil)}
		Services.Set(name, sli)
		if start {
			sli.Start()
		}
	}
}
func RunRouter(ctx context.Context, onExitChannel chan int) (string, *wsrpc.WebsocketRPCConn, error) {
	fmt.Println("Launching Router...")
	internalOnProcessExit := make(chan int, 1)
	routerLogFile, _ := os.OpenFile(LogFilePath("Router"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	routerCmd := exec.Command(RouterFile, Config.UBot.Args()...)
	routerCmd.Stdout = routerLogFile
	routerCmd.Stderr = routerLogFile
	setDeathsig(routerCmd.SysProcAttr)
	err := routerCmd.Start()
	if err != nil {
		return "", nil, fmt.Errorf("cannot start router: %w", err)
	}
	go func() {
		_ = routerCmd.Wait()
		routerLogFile.Close()
		internalOnProcessExit <- 0
	}()
	managerUrl := GetUBotAddr("ws", "/api/manager")
	getTokenUrl := GetUBotAddr("http", "/api/manager/get_token")
	managerToken := ""
	loginParams := url.Values{}
	loginParams.Add("user", Config.UBot.User)
	loginParams.Add("password", Config.UBot.Password)
retryLoop:
	for retryCount := 0; retryCount < 5; retryCount++ {
		timeout := 5 * time.Second
		if retryCount == 0 {
			timeout = 1 * time.Second
		}
		select {
		case <-ctx.Done():
			err = ctx.Err()
			break retryLoop
		case <-time.After(timeout):
			var resp *http.Response
			var req *http.Request
			req, err = http.NewRequestWithContext(ctx,
				http.MethodPost,
				getTokenUrl,
				strings.NewReader(loginParams.Encode()))
			if err != nil {
				continue
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				continue
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				continue
			}
			var tokenBin []byte
			tokenBin, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			managerToken = string(tokenBin)
			break retryLoop
		}
	}
	if managerToken == "" {
		ExitCmd(routerCmd)
		return "", nil, fmt.Errorf("failed to get the manager token: %w", err)
	}
	managerConn, _, err := websocket.DefaultDialer.Dial(managerUrl, nil)
	if err != nil {
		ExitCmd(routerCmd)
		return "", nil, fmt.Errorf("failed to dial to manager: %w", err)
	}
	managerRPCConn := wsrpc.NewWebsocketRPC().Connect(managerConn)
	go func() {
		managerRPCConn.ServeConn()
		managerConn.Close()
		ExitCmd(routerCmd)
		onExitChannel <- 0
	}()
	go func() {
		for {
			select {
			case <-internalOnProcessExit:
				close(internalOnProcessExit)
				return
			case <-ctx.Done():
				ExitCmd(routerCmd)
			}
		}
	}()
	return managerToken, managerRPCConn, nil
}

func DaemonThread(ctx context.Context) {
	var err error
	onRouterExit := make(chan int, 1)
	ManagerToken, ManagerRPCConn, err = RunRouter(ctx, onRouterExit)
	if err != nil {
		log.Fatalf("failed to run router: %v", err)
	}
	for pair := Services.Oldest(); pair != nil; pair = pair.Next() {
		s := pair.Value.(*ServiceLaunchInfo)
		fmt.Printf("Launching %s...\n", s.Service.ID())
		s.Start()
	}
	for {
		select {
		case <-onRouterExit:
			ManagerToken = ""
			ManagerRPCConn = nil
			select {
			case <-ctx.Done():
				return
			default:
			}
			log.Printf("Router is down, retrying in 3 seconds")
			time.Sleep(3 * time.Second)
			for pair := Services.Oldest(); pair != nil; pair = pair.Next() {
				s := pair.Value.(*ServiceLaunchInfo)
				s.Service.Exit()
			}
			ManagerToken, ManagerRPCConn, err = RunRouter(ctx, onRouterExit)
			if err != nil {
				log.Fatal("failed to run router")
			}
		case warppedS := <-OnServiceExit:
			s := warppedS.(*ServiceLaunchInfo)
			if s.Service.Status() == ServiceStopped {
				continue
			}
			if time.Since(s.LaunchAt) > 1*time.Minute {
				s.RetryIndex = 0
			}
			s.RetryIndex++
			if s.RetryIndex > 5 {
				log.Printf("failed to restart %s after 5 attempts within 1 minute", s.Service.ID())
				s.Service.Stop()
				s.RetryIndex = 0
				continue
			}
			go func() {
				log.Printf("%s is down, retrying in 3 seconds", s.Service.ID())
				time.Sleep(3 * time.Second)
				fmt.Printf("Launching %s...\n", s.Service.ID())
				s.Start()
			}()
		}
	}
}
