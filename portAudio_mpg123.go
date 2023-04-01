package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"

	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/gordonklaus/portaudio"
)

func PlayAudioWithMPG123(portAudioDeviceInfo *portaudio.DeviceInfo) {
	// if len(os.Args) < 2 {
	// 	fmt.Println("missing required argument:  input file name")
	// 	return
	// }
	fmt.Println("Playing.  Press Ctrl-C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	ChkErr(err)

	// fileName := os.Args[1]
	ChkErr(decoder.Open("./piano.mp3"))
	defer decoder.Close()

	// get audio format information
	rate, channels, _ := decoder.GetFormat()

	// make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	// portaudio.Initialize()
	// defer portaudio.Terminate()
	out := make([]int16, 8192)

	params := portaudio.HighLatencyParameters(nil, portAudioDeviceInfo)
	params.Input.Channels = 0
	params.Output.Channels = channels
	params.SampleRate = float64(rate)
	params.FramesPerBuffer = len(out)

	stream, err := portaudio.OpenStream(params, &out)
	// stream, err := portaudio.OpenDefaultStream(0, channels, float64(rate), len(out), &out)
	ChkErr(err)
	defer stream.Close()

	ChkErr(stream.Start())
	defer stream.Stop()
	for {
		audio := make([]byte, 2*len(out))
		_, err = decoder.Read(audio)
		if err == mpg123.EOF {
			break
		}
		ChkErr(err)

		ChkErr(binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out))
		ChkErr(stream.Write())
		select {
		case <-sig:
			return
		case <-loop.stop:
			return
		default:
		}
	}
}
