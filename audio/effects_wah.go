package audio

import (
	"math"

	"github.com/gopxl/beep/v2"
)

const (
	wahResonance = 0.2
	wahLfoFreq   = 2.0
)

type wahStreamer struct {
	streamer  beep.Streamer
	lowL      float64
	lowR      float64
	bandL     float64
	bandR     float64
	phase     float64
	phaseStep float64
}

func applyWah(streamer beep.Streamer) beep.Streamer {
	return &wahStreamer{
		streamer:  streamer,
		phaseStep: wahLfoFreq / sampleRate,
	}
}

func (w *wahStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = w.streamer.Stream(samples)
	for i := range n {
		lfoVal := math.Sin(w.phase * 2 * math.Pi)
		cutoff := 950.0 + (600.0 * lfoVal)
		f := 2.0 * math.Sin(math.Pi*cutoff/sampleRate)
		highL := samples[i][0] - w.lowL - (wahResonance * w.bandL)
		highR := samples[i][1] - w.lowR - (wahResonance * w.bandR)
		w.bandL += f * highL
		w.bandR += f * highR
		w.lowL += f * w.bandL
		w.lowR += f * w.bandR
		outL := samples[i][0]*0.1 + w.bandL*2.5
		outR := samples[i][1]*0.1 + w.bandR*2.5
		samples[i][0] = outL / (1.0 + math.Abs(outL))
		samples[i][1] = outR / (1.0 + math.Abs(outR))
		w.phase += w.phaseStep
		if w.phase >= 1.0 {
			w.phase -= 1.0
		}
	}
	return n, ok
}

func (w *wahStreamer) Err() error {
	return w.streamer.Err()
}
