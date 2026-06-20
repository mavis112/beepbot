package audio

import "github.com/gopxl/beep/v2"

type stereo struct {
	streamer beep.Streamer
}

func applyStereo(streamer beep.Streamer) beep.Streamer {
	return &stereo{
		streamer: streamer,
	}
}

func (s *stereo) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = s.streamer.Stream(samples)

	for i := 0; i < n; i++ {
		samples[i][1] = samples[i][0]
	}
	return n, ok
}

func (s *stereo) Err() error {
	return s.streamer.Err()
}
