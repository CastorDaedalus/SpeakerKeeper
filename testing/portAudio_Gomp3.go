package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"os/signal"

	"github.com/gordonklaus/portaudio"
	"github.com/hajimehoshi/go-mp3"
)

func PlayAudioWithGoMP3(outputDevice *portaudio.DeviceInfo) {

	// if len(os.Args) < 2 {
	// 	fmt.Println("missing required argument:  input file name")
	// 	return
	// }
	// fmt.Println("Playing.  Press Ctrl-C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// // create mpg123 decoder instance
	// decoder, err := mpg123.NewDecoder("")
	// chkErr(err)

	// // fileName := os.Args[1]
	// chkErr(decoder.Open("./sound.wav"))
	// defer decoder.Close()

	// // get audio format information
	// rate, channels, _ := decoder.GetFormat()

	// // make sure output format does not change
	// decoder.FormatNone()
	// decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	file, err := os.Open("./piano.mp3")
	ChkErr(err)

	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	// fileBytesReader := bytes.NewReader(fileBytes)

	decoder, err := mp3.NewDecoder(file)
	ChkErr(err)

	out := make([]int16, 8192)

	params := portaudio.HighLatencyParameters(nil, outputDevice)
	params.Input.Channels = 0
	params.Output.Channels = 2
	params.SampleRate = float64(decoder.SampleRate())
	params.FramesPerBuffer = len(out)

	stream, err := portaudio.OpenStream(params, &out)
	// stream, err := portaudio.OpenDefaultStream(0, 2, float64(decoder.SampleRate()), len(out), &out)

	ChkErr(err)
	defer stream.Close()

	ChkErr(stream.Start())
	defer stream.Stop()
	for {
		buffer := make([]byte, 2*len(out))
		buffLength, err := decoder.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		ChkErr(err)

		ChkErr(binary.Read(bytes.NewBuffer(buffer[:buffLength]), binary.LittleEndian, out))
		ChkErr(stream.Write())
		select {
		case <-sig:
			return
		default:
		}
	}
}
