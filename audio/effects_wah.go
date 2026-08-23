package audio

import (
	"math"

	"github.com/gopxl/beep/v2"
)

const (
	// wahResonance = 0.2
	wahLfoFreq = 2.0
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
		currentRes := 0.2 + 0.25*(cutoff-350.0)/1200.0
		f := 2.0 * math.Sin(math.Pi*cutoff/sampleRate)
		highL := samples[i][0] - w.lowL - (currentRes * w.bandL)
		highR := samples[i][1] - w.lowR - (currentRes * w.bandR)
		w.bandL += f * highL
		w.bandR += f * highR
		w.lowL += f * w.bandL
		w.lowR += f * w.bandR
		outL := samples[i][0]*0.45 + w.bandL*1.4
		outR := samples[i][1]*0.45 + w.bandR*1.4
		samples[i][0] = outL / (1.0 + math.Abs(outL)) * 1.2
		samples[i][1] = outR / (1.0 + math.Abs(outR)) * 1.2
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
