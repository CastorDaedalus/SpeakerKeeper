package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

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

func loadConfigData() Config {

	// Check if the file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		println("Creating config file")
		configFile, err := os.Create(filename)
		ChkErr(err)
		defer configFile.Close()
		return baseConfig
	}

	configFile, err := os.Open(filename)
	ChkErr(err)
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&runtimeConfig)

	if err == io.EOF {
		println("Generating base data")
		saveConfigData(runtimeConfig)
		return runtimeConfig
	}
	ChkErr(err)

	fmt.Printf("Loaded config file: %s\n", filename)

	return runtimeConfig

}

func saveConfigData(config Config) {
	configFile, err := os.Create(filename)
	ChkErr(err)
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(config)
	ChkErr(err)

	println("Save complete")
}
