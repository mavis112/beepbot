package bot

import (
	"beepbot/audio"
	"log"
	"sync"
	"sync/atomic"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/gopxl/beep/v2"
)

type Bot struct {
	Client         *twitch.Client
	Channel        string
	player         audio.AudioOutput
	soundsBuffer   map[string]*beep.Buffer
	mtx            sync.RWMutex
	fileMtx        sync.Mutex
	queue          []beep.Streamer
	queueEnabled   bool
	queueIsPlaying bool
	isPlayingSound bool
	speakerIsMuted bool
	erIsOn         bool
	ttsLanguages   map[string]string
	ttsCounter     atomic.Uint64
	volume         int
	playbackID     uint64
}

func New(channel string, player audio.AudioOutput, soundsBuffer map[string]*beep.Buffer, queue bool, er bool, ttsLanguages map[string]string, volume int) *Bot {
	b := &Bot{
		Client:         twitch.NewAnonymousClient(),
		Channel:        channel,
		player:         player,
		soundsBuffer:   soundsBuffer,
		queue:          make([]beep.Streamer, 0, 50),
		queueEnabled:   queue,
		queueIsPlaying: false,
		isPlayingSound: false,
		speakerIsMuted: false,
		erIsOn:         er,
		ttsLanguages:   ttsLanguages,
		volume:         volume,
		playbackID:     0,
	}
	return b
}

func (b *Bot) IsMuted() bool {
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	return b.speakerIsMuted
}

func (b *Bot) IsQueueEnabled() bool {
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	return b.queueEnabled
}

func (b *Bot) SetQEnabled(enabled bool) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.queueEnabled = enabled
}

func (b *Bot) PrintState() {
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	log.Printf("Status -> Mute: %t | Queue: %t | Er: %t | Volume: %d", b.speakerIsMuted, b.queueEnabled, b.erIsOn, b.volume)
}
