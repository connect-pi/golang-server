package main

import (
	"fmt"
	"os"
	"path/filepath"
	"project/pakages/configs"
	"project/pakages/proxy"
	"project/pakages/proxy/rules"
	"project/pakages/v2ray"
	"project/pakages/v2ray/configsTest"
	v2raycore "project/pakages/v2ray/core"
)

func init() {
	// Download v2ray core
	v2raycore.Load()

	// Set core dir
	rootPath, _ := os.Getwd()
	v2ray.CoreDir = filepath.Join(rootPath, "v2ray-core")

	// Create config files
	err := configs.CreateFiles()
	if err != nil {
		fmt.Println("Create config files:", err)
		return
	}

	// Load configs
	err = configs.LoadSettings()
	if err != nil {
		fmt.Println("Load settings:", err)
		return
	}

	err = rules.LoadCustomRules()
	if err != nil {
		fmt.Println("Load settings:", err)
		return
	}
}

func main() {
	// Load Subscription
	if LoadSubscriptionErr := v2ray.LoadSubscription(); LoadSubscriptionErr != nil {
		fmt.Println(LoadSubscriptionErr)
		return
	}

	// Find best config
	bestConfigIndex := configsTest.Run()

	// V2ray Connect
	if connectErr := v2ray.Connect(bestConfigIndex); connectErr != nil {
		fmt.Println(connectErr)
		return
	}

	// Print GO
	fmt.Println(`
░▒▓██████▓▒░ ░▒▓██████▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░
░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
░▒▓█▓▒▒▓███▓▒░▒▓█▓▒░░▒▓█▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░
 ░▒▓██████▓▒░ ░▒▓██████▓▒░
 `)

	// proxy.StartAppHttpProxy(":1080")

	// Create proxy
	proxy.StartAppSocks5Proxy()
}
