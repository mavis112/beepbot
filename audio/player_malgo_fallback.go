//go:build !windows || !cgo

package audio

import (
	"errors"

	"github.com/gopxl/beep/v2"
)

type MalgoPlayer struct{}

func NewMalgoPlayer(sampleRate int, deviceName string) (AudioOutput, error) {
	return nil, errors.New("hello from malgo fallback, switching to oto")
}

func (m *MalgoPlayer) Play(streamer beep.Streamer) {
}

func (m *MalgoPlayer) Stop() {
}

func (m *MalgoPlayer) Close() {
}
