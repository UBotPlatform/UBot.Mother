// +build !dev

package main

import (
	"embed"
	"net/http"
)

//go:generate go run webui_static_generate.go
//go:embed webui/build
var WebUIStaticAssets embed.FS
var WebUIStaticHandler = http.FileServer(http.FS(WebUIStaticAssets))
