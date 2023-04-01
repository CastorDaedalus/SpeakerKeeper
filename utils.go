package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ChkErrMsg(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}

func ChkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
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
