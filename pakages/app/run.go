package app

import (
	"fmt"
	"project/pakages/proxy"
	"project/pakages/v2ray"
	"project/pakages/v2ray/configsTest"
)

func IsRun() bool {
	return v2ray.MainV2RayProcess != nil && v2ray.MainV2RayProcess.IsRun
}

func Start() {
	if IsRun() {
		return
	}

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
	goPint := `
  ░▒▓██████▓▒░ ░▒▓██████▓▒░
  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░
  ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
  ░▒▓█▓▒▒▓███▓▒░▒▓█▓▒░░▒▓█▓▒░
  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░
  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░
   ░▒▓██████▓▒░ ░▒▓██████▓▒░
   `
	fmt.Println(goPint)
	// Logger.Fatalln(goPint)

	// Create proxy
	proxy.StartAppSocks5Proxy()
}

func Stop() {
	v2ray.MainV2RayProcess.Stop(true)
}
