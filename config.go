package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type AccountInfo struct {
	Type string   `yaml:"type,omitempty" json:"type,omitempty"`
	Args []string `yaml:"args,omitempty,flow" json:"args,omitempty"`
}

type UBotConfig struct {
	Address  string `yaml:"address,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
}

func (c *UBotConfig) Args() []string {
	var r []string
	if c.Address != "" {
		r = append(r, "-addr", c.Address)
	}
	if c.User != "" {
		r = append(r, "-user", c.User)
	}
	if c.Password != "" {
		r = append(r, "-password", c.Password)
	}
	return r
}

type WebUIConfig struct {
	Address  string `yaml:"address,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type ConfigInfo struct {
	UBot     UBotConfig    `yaml:"ubot,omitempty"`
	WebUI    WebUIConfig   `yaml:"webui,omitempty"`
	Accounts []AccountInfo `yaml:"accounts,omitempty"`
}

var Config ConfigInfo

func LoadConfig() {
	Config = ConfigInfo{}
	configBinary, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(configBinary, &Config)
	if err != nil {
		log.Println("failed to parse config:", err)
	}
}

func SaveConfig() {
	configBinary, err := yaml.Marshal(Config)
	if err != nil {
		log.Println("failed to marshel config:", err)
		return
	}
	err = ioutil.WriteFile(ConfigFile, configBinary, 0644)
	if err != nil {
		log.Println("failed to write config:", err)
	}
}
