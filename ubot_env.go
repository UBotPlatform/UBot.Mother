package main

import (
	"errors"
	"os"
	"path/filepath"
)

var RootFolder string
var AccountFolder string
var AppFolder string
var SystemFolder string
var ConfigFile string
var LogFolder string
var RouterFile string
var RouterImageName string

func FindUBotEnv() error {
	AccountFolder = filepath.Join(RootFolder, "Accounts")
	AppFolder = filepath.Join(RootFolder, "Apps")
	SystemFolder = filepath.Join(RootFolder, "System")
	LogFolder = filepath.Join(RootFolder, "Logs")
	_ = os.Mkdir(LogFolder, 0755)
	RouterImageName = "Router.ubot" + ExeSuffix
	RouterFile = filepath.Join(SystemFolder, RouterImageName)
	ConfigFile = filepath.Join(RootFolder, "Config.yml")
	if !FileExists(RouterFile) {
		return errors.New("cannot reach the router file (not exist or have no permission)")
	}
	return nil
}
