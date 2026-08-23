package audio

import (
	"github.com/gopxl/beep/v2"
)

type stStreamer struct {
	streamer  beep.Streamer
	buffer    [][2]float64
	chunkSize int
	count     int
	repIndex  int
	readIndex int
	preFilled bool
}

func applyStutter(streamer beep.Streamer, count, intervalMs int) beep.Streamer {
	chunkSize := 44100 * intervalMs / 1000

	return &stStreamer{
		streamer:  streamer,
		buffer:    make([][2]float64, chunkSize),
		chunkSize: chunkSize,
		count:     count,
	}
}

func (s *stStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if !s.preFilled {
		readCount, _ := s.streamer.Stream(s.buffer)
		if readCount < s.chunkSize {
			s.chunkSize = readCount
		}
		s.preFilled = true
	}
	if s.repIndex >= s.count {
		return s.streamer.Stream(samples)
	}
	for i := range len(samples) {
		if s.repIndex >= s.count {
			readCount, ok := s.streamer.Stream(samples[i:])
			return i + readCount, ok
		}
		samples[i][0] = s.buffer[s.readIndex][0]
		samples[i][1] = s.buffer[s.readIndex][1]
		s.readIndex++

		if s.readIndex >= s.chunkSize {
			s.readIndex = 0
			s.repIndex++
		}

	}
	return len(samples), true
}

func (s *stStreamer) Err() error {
	return s.streamer.Err()
}
