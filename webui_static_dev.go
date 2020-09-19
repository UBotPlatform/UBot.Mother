// +build dev

package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var WebUIStaticHandler = func() http.Handler {
	parsedURL, err := url.Parse("http://localhost:3000")
	if err != nil {
		log.Fatalln(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
}()
