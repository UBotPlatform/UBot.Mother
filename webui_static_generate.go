// +build ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	_, sourceFile, _, _ := runtime.Caller(0)
	if _, err := os.Stat(filepath.Join(filepath.Dir(sourceFile), "webui", "node_modules")); err != nil {
		npmCiCmd := exec.Command("npm", "ci")
		npmCiCmd.Dir = filepath.Join(filepath.Dir(sourceFile), "webui")
		err := npmCiCmd.Run()
		if err != nil {
			log.Fatalln(err)
		}
	}
	staticBuildCmd := exec.Command("npm", "run", "build")
	staticBuildCmd.Dir = filepath.Join(filepath.Dir(sourceFile), "webui")
	err := staticBuildCmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
