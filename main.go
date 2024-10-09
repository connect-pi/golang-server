package main

import (
	"project/pakages/api"
	"project/pakages/app"
)

func init() {
	// Start save logs
	app.StartSaveLogs()
}

func main() {
	api.Register()
}
