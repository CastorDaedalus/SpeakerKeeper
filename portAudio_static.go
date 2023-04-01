package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/gordonklaus/portaudio"
)

func playStatic() {
	// portaudio.Initialize()
	// defer portaudio.Terminate()
	h, err := portaudio.DefaultHostApi()
	ChkErr(err)
	stream, err := portaudio.OpenStream(portaudio.HighLatencyParameters(nil, h.DefaultOutputDevice), func(out []int32) {
		for i := range out {
			out[i] = int32(rand.Uint32())
		}
	})
	ChkErr(err)
	defer stream.Close()
	ChkErr(stream.Start())
	time.Sleep(time.Second)
	ChkErr(stream.Stop())
}

const sampleRate = 44100

func playSineWave() {
	portaudio.Initialize()
	defer portaudio.Terminate()
	s := newStereoSine(256, 320, sampleRate)
	defer s.Close()
	ChkErr(s.Start())
	time.Sleep(2 * time.Second)
	ChkErr(s.Stop())
}

type stereoSine struct {
	*portaudio.Stream
	stepL, phaseL float64
	stepR, phaseR float64
}

func newStereoSine(freqL, freqR, sampleRate float64) *stereoSine {
	s := &stereoSine{nil, freqL / sampleRate, 0, freqR / sampleRate, 0}
	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 2, sampleRate, 0, s.processAudio)
	ChkErr(err)
	return s
}

func (g *stereoSine) processAudio(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2 * math.Pi * g.phaseL))
		_, g.phaseL = math.Modf(g.phaseL + g.stepL)
		out[1][i] = float32(math.Sin(2 * math.Pi * g.phaseR))
		_, g.phaseR = math.Modf(g.phaseR + g.stepR)
	}
}
