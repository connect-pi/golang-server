package app

import (
	"fmt"
	"os"
	"path/filepath"
	"project/pakages/configs"
	"project/pakages/proxy"
	"project/pakages/proxy/rules"
	"project/pakages/v2ray"
	v2raycore "project/pakages/v2ray/core"
)

func initApp() {

	// Download v2ray core
	v2raycore.Load()

	// Set core dir
	rootPath, _ := os.Getwd()
	v2ray.CoreDir = filepath.Join(rootPath, ".v2ray-core")

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

func IsRun() bool {
	return v2ray.MainV2RayProcess != nil && v2ray.MainV2RayProcess.IsRun
}

func Start() {
	if IsRun() {
		return
	}

	// initApp()

	// Create proxy
	go proxy.StartAppSocks5Proxy()

	// Load Subscription
	// if LoadSubscriptionErr := v2ray.LoadSubscription(); LoadSubscriptionErr != nil {
	// 	fmt.Println(LoadSubscriptionErr)
	// 	return
	// }

	// // Find best config
	// bestConfigIndex := configsTest.Run()

	// V2ray Connect
	// if connectErr := v2ray.Connect(bestConfigIndex); connectErr != nil {
	// 	fmt.Println(connectErr)
	// 	return
	// }
	if connectErr := v2ray.Connect(0); connectErr != nil {
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

	// Proxy server should keep running
	// for {
	// 	time.Sleep(1 * time.Second)
	// }
	select {}
}

func Stop() {
	v2ray.MainV2RayProcess.Stop(true)
}
