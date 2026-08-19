package audio

import "github.com/gopxl/beep/v2"

const (
	delayBufferSizeL = 8389
	delayBufferSizeR = 9221
	delayFeedback    = 0.45
	tailSamplesMax   = 31000
	delayDryGain     = 0.85
	delayWetGain     = 0.3
)

type delayStreamer struct {
	streamer    beep.Streamer
	lBuffer     []float64
	rBuffer     []float64
	lCounter    int
	rCounter    int
	tailSamples int
}

func applyDelay(streamer beep.Streamer) beep.Streamer {
	return &delayStreamer{
		streamer:    streamer,
		lBuffer:     make([]float64, delayBufferSizeL),
		rBuffer:     make([]float64, delayBufferSizeR),
		lCounter:    0,
		rCounter:    0,
		tailSamples: tailSamplesMax,
	}
}

func (d *delayStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if d.tailSamples <= 0 {
		return 0, false
	}
	n, ok = d.streamer.Stream(samples)
	originalOk := ok
	if !originalOk && d.tailSamples > 0 {
		ok = true
	}
	limit := 0
	if originalOk {
		limit = n
	} else {
		limit = len(samples)
	}
	i := 0

	for ; i < limit; i++ {
		dryL := 0.0
		dryR := 0.0
		if i < n {
			dryL = samples[i][0]
			dryR = samples[i][1]
		} else {
			d.tailSamples--
		}
		echoL := d.lBuffer[d.lCounter]
		echoR := d.rBuffer[d.rCounter]

		outL := dryL*delayDryGain + echoL*delayWetGain
		outR := dryR*delayDryGain + echoR*delayWetGain

		samples[i][0] = outL
		samples[i][1] = outR

		d.rBuffer[d.rCounter] = dryL + (echoL * delayFeedback)
		d.lBuffer[d.lCounter] = dryR + (echoR * delayFeedback)

		d.lCounter++
		if d.lCounter >= len(d.lBuffer) {
			d.lCounter = 0
		}
		d.rCounter++
		if d.rCounter >= len(d.rBuffer) {
			d.rCounter = 0
		}
		if d.tailSamples == 0 {
			i++
			ok = false
			break
		}
	}
	return i, ok
}

func (d *delayStreamer) Err() error {
	return d.streamer.Err()
}
