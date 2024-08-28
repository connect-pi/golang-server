package main

import (
	"project/pakages/configs"
	"project/pakages/proxy"
)

func init() {
	configs.CreateFiles()
}

func main() {

	proxy.Start(":1080")
}
