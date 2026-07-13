package audio

import "github.com/gopxl/beep/v2"

type AudioOutput interface {
	Play(streamer beep.Streamer)
	Stop()
	Close()
}
