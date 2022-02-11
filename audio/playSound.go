package audio

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"log"
	"os"
	"time"
)

const audioFile = "audio/notify.mp3"

func Notify(done chan bool) {
	f, err := os.Open(audioFile)
	if err != nil {
		log.Println(err)
		return
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(streamer beep.StreamSeekCloser) {
		err := streamer.Close()
		if err != nil {
			log.Println(err)
		}
	}(streamer)
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Println(err)
		return
	}
	// done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}
