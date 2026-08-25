package audio

import (
	"errors"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/gopxl/beep/v2"
)

type SoundWithParam struct {
	names           []string
	cutStartPercent int
	cutEndPercent   int
	reversed        bool
	stutterCount    int
	stutterInterval int
	lowQuality      bool
	earRape         bool
	delay           bool
	vibrato         bool
	ringMod         int
	tapeStop        bool
	gacha           bool
	speedRatio      int
}

func CreateSoundWithParam(sounds string, effects string, trackBuffer map[string]*beep.Buffer, isErOn bool) *SoundWithParam {
	namesSlice := strings.Split(sounds, "+")
	soundWithParam := &SoundWithParam{
		names:           []string{},
		cutStartPercent: 0,
		cutEndPercent:   0,
		reversed:        false,
		stutterCount:    0,
		stutterInterval: 0,
		lowQuality:      false,
		earRape:         false,
		delay:           false,
		vibrato:         false,
		ringMod:         0,
		tapeStop:        false,
		gacha:           false,
		speedRatio:      100,
	}
	for _, n := range namesSlice {
		n = strings.ToLower(n)
		if n == "rand" {
			name := getRandomName(trackBuffer)
			soundWithParam.names = append(soundWithParam.names, name)
			continue
		}

		_, ok := trackBuffer[n]
		if !ok {
			continue
		}
		soundWithParam.names = append(soundWithParam.names, n)
	}

	params := strings.Split(effects, "-")

	parseParam(soundWithParam, params)

	if !isErOn {
		soundWithParam.earRape = false
	}

	if soundWithParam.gacha {
		soundWithParam.applyRandomEffects(isErOn)
	}
	return soundWithParam
}

func parseParam(soundWithParam *SoundWithParam, params []string) {
	for _, p := range params {
		if len(p) < 2 {
			continue
		}
		switch string(p[:2]) {
		case "cs":
			cutStartPercent, err := strconv.ParseInt(string(p[2:]), 10, 64)
			if err != nil {
				continue
			}
			if cutStartPercent < 0 {
				cutStartPercent = 1
			}
			if cutStartPercent > 100 {
				cutStartPercent = 100
			}
			soundWithParam.cutStartPercent = int(cutStartPercent)

		case "ce":
			cutEndPercent, err := strconv.ParseInt(string(p[2:]), 10, 64)
			if err != nil {
				continue
			}
			if cutEndPercent < 0 {
				cutEndPercent = 0
			}
			if cutEndPercent > 100 {
				cutEndPercent = 100
			}
			soundWithParam.cutEndPercent = int(cutEndPercent)
		case "rs":
			soundWithParam.reversed = true
		case "st":
			var (
				count    int64 = 3
				interval int64 = 140
				err      error
			)
			if p[2:] == "" {
				soundWithParam.stutterCount = int(count)
				soundWithParam.stutterInterval = int(interval)
				continue
			}
			parts := strings.Split(p[2:], "_")
			count, err = strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				count = 3
			}
			if len(parts) > 1 {

				interval, err = strconv.ParseInt(parts[1], 10, 64)
				if err != nil {
					interval = 140
				}
			}

			if count < 1 {
				count = 1
			}
			if count > 8 {
				count = 8
			}
			if interval < 60 {
				interval = 60
			}
			if interval > 300 {
				interval = 300
			}
			soundWithParam.stutterCount = int(count)
			soundWithParam.stutterInterval = int(interval)

		case "lq":
			soundWithParam.lowQuality = true
		case "er":
			soundWithParam.earRape = true
		case "dl":
			soundWithParam.delay = true
		case "vb":
			soundWithParam.vibrato = true

		case "rm":
			ringMod, err := strconv.ParseInt(string(p[2:]), 10, 64)
			if err != nil {
				ringMod = 50
			}
			if ringMod < 0 {
				ringMod = 1
			}
			if ringMod > 100 {
				ringMod = 100
			}
			soundWithParam.ringMod = int(ringMod)
		case "ts":
			soundWithParam.tapeStop = true
		case "ga":
			soundWithParam.gacha = true
		case "sp":
			speedRatio, err := strconv.ParseInt(string(p[2:]), 10, 64)
			if err != nil {
				continue
			}
			if speedRatio < 10 {
				speedRatio = 10
			}
			if speedRatio > 200 {
				speedRatio = 200
			}
			soundWithParam.speedRatio = int(speedRatio)
		}
	}
}

func (s *SoundWithParam) applyRandomEffects(isErOn bool) {
	appliedC := 0
	candidates := make([]string, 0, 9)
	if s.reversed {
		appliedC++
	} else {
		candidates = append(candidates, "reversed")
	}
	if s.stutterCount != 0 {
		appliedC++
	} else {
		candidates = append(candidates, "stutter")
	}
	if s.lowQuality {
		appliedC++
	} else {
		candidates = append(candidates, "lowQuality")
	}
	if s.earRape {
		appliedC++
	} else {
		if isErOn {
			candidates = append(candidates, "earRape")
		}
	}
	if s.delay {
		appliedC++
	} else {
		candidates = append(candidates, "delay")
	}
	if s.vibrato {
		appliedC++
	} else {
		candidates = append(candidates, "vibrato")
	}
	if s.ringMod != 0 {
		appliedC++
	} else {
		candidates = append(candidates, "ringMod")
	}
	if s.tapeStop {
		appliedC++
	} else {
		candidates = append(candidates, "tapeStop")
	}

	if s.speedRatio != 100 {
		appliedC++
	} else {
		candidates = append(candidates, "speed")
	}

	limit := 3 - appliedC
	if limit < 0 {
		limit = 0
	}

	limit = isCrit(limit)

	if limit == 0 {
		return
	}

	finalCount := getFinalCount(limit)

	rand.Shuffle(len(candidates), func(i int, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	for i := 0; i < finalCount; i++ {
		switch candidates[i] {
		case "reversed":
			s.reversed = true
		case "stutter":
			s.stutterCount = rand.IntN(8) + 1
			s.stutterInterval = rand.IntN(241) + 60
		case "lowQuality":
			s.lowQuality = true
		case "earRape":
			s.earRape = true
		case "delay":
			s.delay = true
		case "vibrato":
			s.vibrato = true
		case "ringMod":
			s.ringMod = rand.IntN(91) + 10
		case "ts":
			s.tapeStop = true
		case "speed":
			s.speedRatio = randomSpeedRatio()
		}
	}
}

func isCrit(num int) int {
	roll := rand.IntN(100) + 1
	if roll <= 5 {
		return num + 1
	}
	return num
}

func getFinalCount(n int) int {
	switch n {
	case 1:
		return 1
	case 2:
		roll := rand.IntN(100) + 1
		if roll <= 70 {
			return 1
		} else {
			return 2
		}
	case 3:
		roll := rand.IntN(100) + 1
		if roll <= 50 {
			return 1
		} else if roll <= 85 {
			return 2
		} else {
			return 3
		}
	case 4:
		roll := rand.IntN(100) + 1
		if roll <= 30 {
			return 2
		} else if roll <= 80 {
			return 3
		} else {
			return 4
		}
	}
	return 0
}

func randomSpeedRatio() int {
	speed := 0
	roll := rand.IntN(100) + 1

	if roll <= 45 {
		speed = rand.IntN(80-50+1) + 50
	} else if roll <= 85 {
		speed = rand.IntN(170-120+1) + 120
	} else {
		speed = rand.IntN(45-20+1) + 20
	}
	return speed
}

func CreateStreamerWithParameter(s *SoundWithParam, trackBuffer map[string]*beep.Buffer) (beep.Streamer, error) {
	if len(s.names) < 1 {
		return nil, errors.New("audio is empty")
	}
	var totalLen int
	streamerSlice := make([]beep.Streamer, 0, len(s.names))
	var streamer beep.Streamer
	var maxLen int
	for _, name := range s.names {
		currBuffer := trackBuffer[name]
		if currBuffer == nil {
			return nil, errors.New("corrupted or empty buffer")
		}
		totalLen = currBuffer.Len()
		start := totalLen * s.cutStartPercent / 100
		end := totalLen * (100 - s.cutEndPercent) / 100
		if start >= end {
			start = 0
			end = totalLen
		}
		if currLen := end - start; currLen > maxLen {
			maxLen = currLen
		}
		var str beep.Streamer
		if s.reversed {
			if end-start > maxReverseSamples {
				end = start + maxReverseSamples
			}
			str = applyReverse(currBuffer, start, end)
		} else {
			str = currBuffer.Streamer(start, end)
		}
		streamerSlice = append(streamerSlice, str)
	}
	streamer = beep.Mix(streamerSlice...)
	if s.tapeStop {
		streamer = applyTs(streamer, maxLen)
	}
	if s.stutterCount != 0 {
		streamer = applyStutter(streamer, s.stutterCount, s.stutterInterval)
	}
	if s.lowQuality {
		streamer = applyLowQuality(streamer)
	}
	if s.earRape {
		streamer = applyEarRape(streamer)
	}
	if s.vibrato {
		streamer = applyVibrato(streamer)
	}
	if s.ringMod != 0 {
		streamer = applyRingMod(streamer, s.ringMod)
	}
	if s.delay {
		streamer = applyDelay(streamer)
	}
	if s.speedRatio != 100 {
		streamer = beep.ResampleRatio(3, float64(s.speedRatio)/100.0, streamer)
	}

	return streamer, nil
}
