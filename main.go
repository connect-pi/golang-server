package main

import (
	"os"
	"path/filepath"
	"project/pakages/api"
	"project/pakages/clog"
	"project/pakages/configs"
	"project/pakages/proxy"
	"project/pakages/proxy/rules"
	"project/pakages/v2ray"
	v2raycore "project/pakages/v2ray/core"
)

func init() {
	// Set core dir
	rootPath, _ := os.Getwd()
	v2ray.CoreDir = filepath.Join(rootPath, ".v2ray-core")

	// Create config files
	err := configs.CreateFiles()
	if err != nil {
		clog.Println("Create config files:", err)
		return
	}

	// Load custom rules
	err = rules.LoadCustomRules()
	if err != nil {
		clog.Println("Load settings:", err)
		return
	}

	// Load settings
	err = configs.LoadSettings()
	if err != nil {
		clog.Println("Load settings:", err)
		return
	}

	// Download v2ray core
	v2raycore.Load()
}

func main() {
	// Create proxy
	go proxy.StartAppSocks5Proxy()

	api.Register()
}
