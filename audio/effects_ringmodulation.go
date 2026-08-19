package audio

import (
	"math"

	"github.com/gopxl/beep/v2"
)

const (
	ringModFreqMult = 10.0
)

type ringModStreamer struct {
	streamer  beep.Streamer
	phase     float64
	phaseStep float64
}

func applyRingMod(streamer beep.Streamer, level int) beep.Streamer {
	freq := ringModFreqMult * float64(level)
	return &ringModStreamer{
		streamer:  streamer,
		phase:     0,
		phaseStep: freq / sampleRate,
	}
}

func (r *ringModStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = r.streamer.Stream(samples)
	for i := range n {
		carrier := math.Sin(r.phase * 2.0 * math.Pi)
		samples[i][0] *= carrier
		samples[i][1] *= carrier
		r.phase += r.phaseStep
		if r.phase >= 1.0 {
			r.phase -= 1.0
		}
	}
	return n, ok
}

func (r *ringModStreamer) Err() error {
	return r.streamer.Err()
}
