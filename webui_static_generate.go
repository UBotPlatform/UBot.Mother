// +build ignore

package main

import (
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/shurcooL/vfsgen"
)

func main() {
	_, sourceFile, _, _ := runtime.Caller(0)
	staticBuildCmd := exec.Command("npm", "run", "build")
	staticBuildCmd.Dir = filepath.Join(filepath.Dir(sourceFile), "webui")
	err := staticBuildCmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	err = vfsgen.Generate(http.Dir("webui/build"), vfsgen.Options{
		PackageName:  "main",
		BuildTags:    "!dev",
		VariableName: "WebUIStaticAssets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
