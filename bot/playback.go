package bot

import (
	"beepbot/audio"
	"math"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
)

func (b *Bot) assembleStreamer(taskSlice []PlayTask) beep.Streamer {
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	streamersSlice := make([]beep.Streamer, 0, len(taskSlice))
	for _, t := range taskSlice {
		sound := audio.CreateSoundWithParam(t.Content, t.Effects, b.soundsBuffer, b.erIsOn)
		streamer, err := audio.CreateStreamerWithParameter(sound, b.soundsBuffer)
		if err != nil {
			continue
		}
		streamersSlice = append(streamersSlice, streamer)
	}

	if len(streamersSlice) < 1 {
		return nil
	}
	streamer := beep.Seq(streamersSlice...)

	if b.volume != 100 {
		vol, silent := getVol(b.volume)
		streamer = &effects.Volume{
			Streamer: streamer,
			Base:     2,
			Volume:   vol,
			Silent:   silent,
		}
	}
	return streamer
}

func getVol(vol int) (float64, bool) {
	if vol == 0 {
		return 0, true
	}
	return math.Log2(float64(vol) / 100), false
}

func (b *Bot) pushToQueue(s beep.Streamer) {
	b.mtx.Lock()
	if len(b.queue) >= 50 {
		b.mtx.Unlock()
		return
	}
	b.queue = append(b.queue, s)
	if !b.queueIsPlaying {
		b.queueIsPlaying = true
		b.mtx.Unlock()
		b.playNext()
		return
	}
	b.mtx.Unlock()
}

func (b *Bot) playNext() {
	b.mtx.Lock()
	if len(b.queue) == 0 {
		b.queueIsPlaying = false
		b.mtx.Unlock()
		return
	}
	nextSound := b.queue[0]

	b.queue[0] = nil

	b.queue = b.queue[1:]

	b.isPlayingSound = true

	b.playbackID++
	currID := b.playbackID
	b.mtx.Unlock()
	b.player.Play(beep.Seq(nextSound, beep.Callback(func() {
		b.mtx.Lock()
		if currID != b.playbackID {
			b.mtx.Unlock()
			return
		}
		b.isPlayingSound = false
		b.mtx.Unlock()
		go b.playNext()
	})))
}

func (b *Bot) deleteKeys(keys []string) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	for _, key := range keys {
		delete(b.soundsBuffer, key)
	}
}
