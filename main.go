package main

import (
	"fmt"
	"project/pakages/configs"
	"project/pakages/proxy"
	"project/pakages/v2ray"
)

func init() {
	// Config files
	configs.CreateFiles()
	configs.LoadSettings()
}

func main() {
	// Load Subscription
	if configs.Settings.UpdateSubscription {
		if LoadSubscriptionErr := v2ray.LoadSubscription(); LoadSubscriptionErr != nil {
			fmt.Println(LoadSubscriptionErr)
			return
		}
	}

	// V2ray Connect
	if connectErr := v2ray.Connect(0); connectErr != nil {
		fmt.Println(connectErr)
		return
	}

	// Create proxy
	proxy.Start(":1080")
}
