package sound

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/faiface/beep/speaker"
)

func PlaySound(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening sound file: %v", err)
		return
	}
	defer file.Close()

	streamer, format, err := wav.Decode(file)
	if err != nil {
		log.Printf("Error decoding sound file: %v", err)
		return
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Printf("Error initializing speaker: %v", err)
		return
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
	log.Println("Sound played successfully")
}
