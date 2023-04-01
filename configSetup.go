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

	if runtimeConfig.SelectedDeviceIndex < 0 {
		outputDevicesStartIndex, err := getOutputDeviceIndex(portAudioDevices)
		ChkErr(err)
		guiSelectedDeviceName, err = getUserSelectedAudioDevice(portAudioDevices[outputDevicesStartIndex:])
		ChkErr(err)

		for i, portAudioDevice := range portAudioDevices {
			if guiSelectedDeviceName == portAudioDevice.Name {
				runtimeConfig.SelectedDeviceIndex = i
				needToSave = true
			}
		}

		if runtimeConfig.SelectedDeviceIndex < 0 {
			err = errors.New("unable to find audio device")
			ChkErr(err)
		}

	}

	if runtimeConfig.TimeInterval <= 0 {
		runtimeConfig.TimeInterval = getUserInputWaitTime()
		needToSave = true
	}

	if needToSave {
		saveConfigData(runtimeConfig)
	}

}

func getOutputDeviceIndex(portAudioDevices []*portaudio.DeviceInfo) (int, error) {

	var err error

	for i := 0; i < len(portAudioDevices); i++ {
		if strings.Contains(portAudioDevices[i].Name, "Microsoft Sound Mapper - Output") ||
			strings.Contains(portAudioDevices[i].Name, "Default Output Device") {
			portAudioDevices[i].Name = "Default Output Device"
			return i, err
		}
	}

	return -1, errors.New("couldn't find primary output device in portaudio devices")

}
