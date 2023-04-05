package main

import (
	"fmt"
	"math"
	"speakerkeeper/icon"
	"strconv"
	"time"

	"github.com/getlantern/systray"
	"github.com/gordonklaus/portaudio"
	"github.com/ncruces/zenity"
)

type audioLoop struct {
	ticker *time.Ticker
	stop   chan struct{}
	active bool
}

var loop = audioLoop{}

func initMainLoop() {
	onExit := func() {
		println("Done")
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	if isConfigValid() {
		loop.startAudioLoop()
	}

	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("SpeakerKeeper")
	systray.SetTooltip("SpeekerKeeper")

	reset := systray.AddMenuItem("Reset Config", "Erase data and start setup again")
	stop := systray.AddMenuItem("Stop Playing", "Stop playing audio file")
	start := systray.AddMenuItem("Start Playing", "Start playing audio file")
	systray.AddSeparator()
	quitBtn := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		if loop.active {
			start.Hide()
		} else {
			stop.Hide()
		}

		playing := true
		toggle := func() {
			if playing {
				start.Show()
				stop.Hide()
				playing = false
			} else {
				start.Hide()
				stop.Show()
				playing = true
			}
		}
		for {
			select {
			case <-reset.ClickedCh:
				loop.stopAudioLoop()
				runtimeConfig = baseConfig
				getNecessaryConfig()
				if isConfigValid() {
					loop.startAudioLoop()
				}
				if playing != loop.active {
					toggle()
				}

			case <-stop.ClickedCh:
				loop.stopAudioLoop()
				toggle()

			case <-start.ClickedCh:
				if isConfigValid() {
					loop.startAudioLoop()
					toggle()
				}
			}
		}
	}()

	go func() {
		<-quitBtn.ClickedCh
		fmt.Println("Requesting quit")
		loop.stopAudioLoop()
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

}

func (loop *audioLoop) startAudioLoop() {
	duration := time.Duration(runtimeConfig.TimeInterval * float64(time.Minute))
	loop.ticker = time.NewTicker(duration)
	loop.stop = make(chan struct{})
	loop.active = true

	go func() {
		fmt.Println("Starting audio loop")
		defer loop.ticker.Stop()
		for {
			select {
			case <-loop.ticker.C:
				go func() {
					devices, err := portaudio.Devices()
					ChkErr(err)
					for _, device := range devices {
						if runtimeConfig.SelectedDeviceName == device.Name {
							PlayAudioWithMPG123(device)
							break
						}
					}
				}()
			case <-loop.stop:
				fmt.Println("Stopping audio loop")
				return
			}
		}
	}()
}

func (loop *audioLoop) stopAudioLoop() {
	if loop.active {
		loop.active = false
		close(loop.stop)
	}
}

func getUserSelectedAudioDevice(portAudioDevices []*portaudio.DeviceInfo) (string, error) {

	fullNameSlice := make([]string, len(portAudioDevices))
	for i, device := range portAudioDevices {
		fullNameSlice[i] = device.Name
	}

	selected, err := zenity.List(
		"Detected Audio Output Devices",
		fullNameSlice,
		zenity.Title("Select an audio device"),
		zenity.Width(300),
		zenity.Height(300),
		zenity.DisallowEmpty(),
		zenity.RadioList(),
		zenity.WindowIcon(zenity.QuestionIcon),
	)

	return selected, err

}

func getUserInputWaitTime() (float64, error) {
	var minutesFloat float64
	var err error
	for {
		minutes, err := zenity.Entry(
			"Time interval in minutes",
			zenity.Title("Sound playing frequency"),
			zenity.WindowIcon(zenity.QuestionIcon),
		)
		if err != nil {
			return 0, err
		}
		minutesFloat, err = strconv.ParseFloat(minutes, 64)
		if err != nil || minutesFloat < 0 {
			zenity.Warning("Please enter only numbers greater than 0",
				zenity.Title("Invalid format"),
				zenity.WindowIcon(zenity.WarningIcon))
		} else {
			break
		}
	}

	return math.Floor(minutesFloat*100) / 100, err
}

func printPortAudioDevicesWithFullNames(portAudioDevices []*portaudio.DeviceInfo) {

	fmt.Println("\nDevice list\n-----------")
	for i, device := range portAudioDevices {
		fmt.Printf("[%d] %s\n", i, device.Name)
	}

}

func getUserInputCMD() int {
	// Ask the user to select a device to play audio from
	var deviceId int
	fmt.Print("\nEnter the ID of the device to play audio from: ")
	fmt.Scanln(&deviceId)

	return deviceId
}
