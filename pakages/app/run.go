package app

import (
	"project/pakages/clog"
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
		clog.Println(LoadSubscriptionErr)
		return
	}

	// Find best config
	bestConfigIndex := configsTest.Run()

	// V2ray Connect
	if connectErr := v2ray.Connect(bestConfigIndex); connectErr != nil {
		clog.Println(connectErr)
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
	clog.Println(goPint)
	// Logger.Fatalln(goPint)

	select {}
}

func Stop() {
	v2ray.MainV2RayProcess.Stop(true)
}
