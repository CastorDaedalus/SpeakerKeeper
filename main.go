package main

import (
	"github.com/gordonklaus/portaudio"
)

var filename = "config.json"
var baseConfig Config = Config{SelectedDeviceIndex: -1}
var runtimeConfig Config = Config{SelectedDeviceIndex: -1}

func main() {

	portaudio.Initialize()
	defer portaudio.Terminate()

	runtimeConfig = loadConfigData()
	getNecessaryConfig()

	initMainLoop()

}
