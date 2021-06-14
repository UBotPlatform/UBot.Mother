// +build !dev

package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:generate go run webui_static_generate.go
//go:embed webui/build/*
var WebUIStaticAssets embed.FS
var WebUIStaticHandler = http.FileServer(http.FS(func() fs.FS {
	r, _ := fs.Sub(WebUIStaticAssets, "webui/build")
	return r
}()))
