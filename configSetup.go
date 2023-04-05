package main

import (
	"errors"
	"strings"

	"github.com/gordonklaus/portaudio"
)

func getNecessaryConfig() {
	var guiSelectedDeviceName string
	portAudioDevices, err := portaudio.Devices()
	ChkErr(err)

	needToSave := false

	if runtimeConfig.PortAudioDeviceList == nil {
		namesList := make([]string, 0, 4)
		for _, device := range portAudioDevices {
			namesList = append(namesList, device.Name)
		}
		runtimeConfig.PortAudioDeviceList = namesList
		needToSave = true
	}

	if runtimeConfig.SelectedDeviceName == "" {
		outputDevicesStartIndex, err := getOutputDeviceIndex(portAudioDevices)
		ChkErr(err)
		guiSelectedDeviceName, err = getUserSelectedAudioDevice(portAudioDevices[outputDevicesStartIndex:])
		if err == nil {
			for _, portAudioDevice := range portAudioDevices {
				if guiSelectedDeviceName == portAudioDevice.Name {
					runtimeConfig.SelectedDeviceName = portAudioDevice.Name
					needToSave = true
				}
			}

			if runtimeConfig.SelectedDeviceName == "" {
				err = errors.New("unable to find audio device")
				panic(err)
			}
		} else {
			return
		}

	}

	if runtimeConfig.TimeInterval <= 0 {
		runtimeConfig.TimeInterval, err = getUserInputWaitTime()
		if err == nil {
			needToSave = true
		}
	}

	if needToSave {
		saveConfigData(runtimeConfig)
	}

}

func getOutputDeviceIndex(portAudioDevices []*portaudio.DeviceInfo) (int, error) {

	var err error

	for i := 0; i < len(portAudioDevices); i++ {
		if strings.Contains(portAudioDevices[i].Name, "Microsoft Sound Mapper - Output") {
			return i, err
		}
	}

	return -1, errors.New("couldn't find primary output device in portaudio devices")

}

func isConfigValid() bool {
	if runtimeConfig.SelectedDeviceName != "" && runtimeConfig.TimeInterval > 0 {
		return true
	} else {
		return false
	}
}
