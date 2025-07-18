package sound

import (
	"math"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

const sampleRate = 44100

var Mute bool = false

type note struct {
	freq     float64
	duration time.Duration
}

func generateTone(frequency float64, duration time.Duration) beep.Streamer {
	streamer := beep.StreamerFunc(func(samples [][2]float64) (int, bool) {
		for i := range samples {
			t := float64(i) / float64(sampleRate)
			v := 0.2 * math.Sin(2*math.Pi*frequency*t)
			samples[i][0] = v
			samples[i][1] = v
		}
		return len(samples), true
	})
	return beep.Take(beep.SampleRate(sampleRate).N(duration), streamer)
}

func playSequence(notes []note) {
	if Mute {
		return 
	}
	speaker.Init(beep.SampleRate(sampleRate), sampleRate/10)

	var streamers []beep.Streamer
	for _, n := range notes {
		tone := generateTone(n.freq, n.duration)
		silence := beep.Silence(beep.SampleRate(sampleRate).N(30 * time.Millisecond))
		streamers = append(streamers, tone, silence)
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(
		append(streamers, beep.Callback(func() { done <- true }))...,
	))
	<-done
}
