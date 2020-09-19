// +build !dev

package main

import (
	"net/http"
)

//go:generate go run webui_static_generate.go
var WebUIStaticHandler = http.FileServer(WebUIStaticAssets)
