package main

import (
	"fmt"
	"project/pakages/configs"
	"project/pakages/proxy"
	"project/pakages/proxy/rules"
	"project/pakages/v2ray"
)

func init() {
	err := configs.CreateFiles()
	if err != nil {
		fmt.Println("Create config files:", err)
		return
	}

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

	// V2ray Connect
	if connectErr := v2ray.Connect(0); connectErr != nil {
		fmt.Println(connectErr)
		return
	}

	// Create proxy
	proxy.Start(":1080")
	// configsTest.CreateJsonFiles()
}
