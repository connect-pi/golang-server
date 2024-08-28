package main

import (
	"project/pakages/configs"
	"project/pakages/proxy"
)

func init() {
	// Create config files
	configs.CreateFiles()
}

func main() {
	// Load Subscription
	// if LoadSubscriptionErr := v2ray.LoadSubscription(); LoadSubscriptionErr != nil {
	// 	fmt.Println(LoadSubscriptionErr)
	// 	return
	// }

	// // V2ray Connect
	// if connectErr := v2ray.Connect(0); connectErr != nil {
	// 	fmt.Println(connectErr)
	// 	return
	// }

	// Create proxy
	proxy.Start(":1080")
}
