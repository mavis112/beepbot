package audio

import "github.com/gopxl/beep/v2"

type tsStreamer struct {
	streamer         beep.Streamer
	buffer           [32768][2]float64
	writeCount       int
	pos              float64
	active           bool
	stopCounter      int
	actualStopWindow int
	triggerPoint     float64
	exhausted        bool
}

func applyTs(streamer beep.Streamer, totalSamples int) beep.Streamer {
	window := 44100
	if totalSamples < 44100 {
		window = totalSamples / 2
	}
	triggerPoint := float64(totalSamples - window)
	return &tsStreamer{
		streamer:         streamer,
		buffer:           [32768][2]float64{},
		actualStopWindow: window,
		triggerPoint:     triggerPoint,
	}
}

func (t *tsStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if !t.exhausted {
		temp := make([][2]float64, len(samples))
		readCount, _ := t.streamer.Stream(temp)
		for _, sample := range temp[:readCount] {
			t.buffer[t.writeCount&32767][0] = sample[0]
			t.buffer[t.writeCount&32767][1] = sample[1]
			t.writeCount++
		}
		if readCount < len(samples) {
			t.exhausted = true
		}
	}
	step := 1.0
	for i := range len(samples) {
		if !t.active && t.pos > t.triggerPoint {
			t.active = true
		}
		if t.active {
			t.stopCounter++
			if t.stopCounter >= t.actualStopWindow || t.pos >= float64(t.writeCount) {
				return i, false
			}
			progress := float64(t.stopCounter) / float64(t.actualStopWindow)
			step = 1 - progress*progress

		}
		index1 := int(t.pos)
		index2 := index1 + 1
		wrapIndex1 := index1 & 32767
		wrapIndex2 := index2 & 32767
		frac := t.pos - float64(index1)
		samples[i][0] = t.buffer[wrapIndex1][0]*(1.0-frac) + t.buffer[wrapIndex2][0]*frac
		samples[i][1] = t.buffer[wrapIndex1][1]*(1.0-frac) + t.buffer[wrapIndex2][1]*frac

		t.pos += step
	}
	return len(samples), true
}

func (t *tsStreamer) Err() error {
	return t.streamer.Err()
}
