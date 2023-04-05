package main

import (
	"github.com/gordonklaus/portaudio"
)

var filename = "config.json"
var baseConfig Config = Config{}
var runtimeConfig Config = Config{}

func main() {

	portaudio.Initialize()
	defer portaudio.Terminate()

	runtimeConfig = loadConfigData()
	getNecessaryConfig()

	initMainLoop()

}
