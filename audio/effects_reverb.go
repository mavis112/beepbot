package audio

import (
	"math"

	"github.com/gopxl/beep/v2"
)

const (
	reverbTailSamples     = 4 * 44100
	reverbDecay           = 0.7995*0.8 + 0.005
	reverbTankFeedback    = 0.5
	reverbFilterFreq      = 0.9*18400 + 100
	reverbFilterStartFreq = 1000.0
	reverbControlRate     = 44
	reverbOversample      = 4
	reverbOutputGain      = 0.6
)

var (
	reverbDiffuserLengths  = [4]int{211, 158, 560, 410}
	reverbDiffuserFeedback = [4]float64{0.75, 0.75, 0.625, 0.625}
)

type reverbStreamer struct {
	streamer    beep.Streamer
	drained     bool
	tailSamples int

	bandwidth      [2]reverbLowPass
	damping        [2]reverbLowPass
	controlCounter int
	predelay       float64
	diffusers      [4]*reverbLine
	tankAllPass    [4]*reverbLine
	tankDelay      [4]*reverbLine
	prevLeftTank   float64
	prevRightTank  float64
}

func applyReverb(streamer beep.Streamer) beep.Streamer {
	r := &reverbStreamer{
		streamer:    streamer,
		tailSamples: reverbTailSamples,
		tankAllPass: [4]*reverbLine{
			newReverbLine(705),
			newReverbLine(2116, 211, 1446, 0),
			newReverbLine(1058),
			newReverbLine(3139, 1093, 388, 0),
		},
		tankDelay: [4]*reverbLine{
			newReverbLine(5292, 2363, 388, 4268),
			newReverbLine(4233, 1270, 3139, 0),
			newReverbLine(4939, 313, 3492, 0),
			newReverbLine(3880, 2363, 144, 0),
		},
	}
	for i, length := range reverbDiffuserLengths {
		r.diffusers[i] = newReverbLine(length)
	}
	for i := range 2 {
		r.bandwidth[i].setFrequency(reverbFilterStartFreq)
		r.damping[i].setFrequency(reverbFilterStartFreq)
	}
	return r
}

func (r *reverbStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if r.tailSamples <= 0 {
		return 0, false
	}
	if !r.drained {
		n, ok = r.streamer.Stream(samples)
		if !ok {
			r.drained = true
		}
	}
	if r.drained {
		for n < len(samples) && r.tailSamples > 0 {
			samples[n] = [2]float64{}
			n++
			r.tailSamples--
		}
	}
	for i := range n {
		samples[i][0], samples[i][1] = r.process(samples[i][0], samples[i][1])
	}
	return n, n > 0
}

func (r *reverbStreamer) Err() error {
	return r.streamer.Err()
}

func (r *reverbStreamer) process(left, right float64) (float64, float64) {
	if r.controlCounter >= reverbControlRate {
		r.controlCounter = 0
		for i := range 2 {
			r.bandwidth[i].setFrequency(reverbFilterFreq)
			r.damping[i].setFrequency(reverbFilterFreq)
		}
	}
	r.controlCounter++

	bandwidthL := r.bandwidth[0].process(left)
	bandwidthR := r.bandwidth[1].process(right)

	input := r.predelay
	r.predelay = (bandwidthL + bandwidthR) * 0.5

	for i, d := range r.diffusers {
		input = d.allPass(input, reverbDiffuserFeedback[i])
	}

	leftTank := r.tankAllPass[0].delay(input + r.prevRightTank)
	leftTank = r.tankDelay[0].delay(leftTank)
	leftTank = r.damping[0].process(leftTank)
	leftTank = r.tankAllPass[1].allPass(leftTank, reverbTankFeedback)
	leftTank = r.tankDelay[1].delay(leftTank)

	rightTank := r.tankAllPass[2].delay(input + r.prevLeftTank)
	rightTank = r.tankDelay[2].delay(rightTank)
	rightTank = r.damping[1].process(rightTank)
	rightTank = r.tankAllPass[3].allPass(rightTank, reverbTankFeedback)
	rightTank = r.tankDelay[3].delay(rightTank)

	r.prevLeftTank = math.Tanh(leftTank * reverbDecay)
	r.prevRightTank = math.Tanh(rightTank * reverbDecay)

	outL := r.tankDelay[2].tap(1) + r.tankDelay[2].tap(2) - r.tankAllPass[3].tap(1) +
		r.tankDelay[3].tap(1) - r.tankDelay[0].tap(1) - r.tankAllPass[1].tap(1) - r.tankDelay[1].tap(1)
	outR := r.tankDelay[0].tap(2) + r.tankDelay[0].tap(3) - r.tankAllPass[1].tap(2) +
		r.tankDelay[1].tap(2) - r.tankDelay[2].tap(3) - r.tankAllPass[3].tap(2) - r.tankDelay[3].tap(2)

	finalL := math.Tanh(outL * reverbOutputGain)
	finalR := math.Tanh(outR * reverbOutputGain)

	return finalL, finalR
}

type reverbLine struct {
	buffer []float64
	taps   []int
}

func newReverbLine(length int, extraTaps ...int) *reverbLine {
	return &reverbLine{
		buffer: make([]float64, length),
		taps:   append([]int{0}, extraTaps...),
	}
}

func (l *reverbLine) delay(input float64) float64 {
	output := l.buffer[l.taps[0]]
	l.buffer[l.taps[0]] = input
	l.advance()
	return output
}

func (l *reverbLine) allPass(input, feedback float64) float64 {
	bufOut := l.buffer[l.taps[0]]
	output := bufOut - input*feedback
	l.buffer[l.taps[0]] = input + output*feedback
	l.advance()
	return output
}

func (l *reverbLine) tap(i int) float64 {
	return l.buffer[l.taps[i]]
}

func (l *reverbLine) advance() {
	for i := range l.taps {
		l.taps[i]++
		if l.taps[i] >= len(l.buffer) {
			l.taps[i] = 0
		}
	}
}

type reverbLowPass struct {
	f    float64
	low  float64
	band float64
}

func (p *reverbLowPass) setFrequency(freq float64) {
	p.f = 2 * math.Sin(math.Pi*freq/(sampleRate*reverbOversample))
}

func (p *reverbLowPass) process(input float64) float64 {
	for range reverbOversample {
		p.low += p.f*p.band + 1e-25
		high := input - p.low - 2*p.band
		p.band += p.f * high
	}
	return p.low
}
