package audio

import (
	"fmt"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

type OtoPlayer struct{}

func NewOtoPlayer(sampleRate int) (AudioOutput, error) {
	sr := beep.SampleRate(sampleRate)
	if err := speaker.Init(sr, sr.N(time.Second/10)); err != nil {
		return nil, fmt.Errorf("player is failed to init: %w", err)
	}
	return &OtoPlayer{}, nil
}

func (o *OtoPlayer) Play(streamer beep.Streamer) {
	speaker.Play(streamer)
}

func (o *OtoPlayer) Stop() {
	speaker.Clear()
}

func (o *OtoPlayer) Close() {
}
