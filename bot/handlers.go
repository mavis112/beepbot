package bot

import (
	"log"
	"strconv"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/joho/godotenv"
)

func (b *Bot) handleMessage(msg twitch.PrivateMessage) {
	if len(msg.Message) < 2 || (!strings.HasPrefix(msg.Message, "!") && !strings.HasPrefix(msg.Message, "@")) {
		return
	}

	msgSlice := strings.Fields(msg.Message)
	if len(msgSlice) == 0 {
		return
	}

	if len(msgSlice) > 2 && strings.HasPrefix(msgSlice[0], "@") {
		msgSlice = msgSlice[1:]
		msg.Message = strings.Join(msgSlice, " ")
	}

	command := strings.ToLower(msgSlice[0])

	if command == "!m" {
		b.playSound(msg)
	}
}

func (b *Bot) HandleLoop(msgChan <-chan twitch.PrivateMessage) {
	for msg := range msgChan {
		b.handleMessage(msg)
	}
}

func (b *Bot) playSound(msg twitch.PrivateMessage) {
	msgSlice := strings.Fields(msg.Message)[1:]
	if len(msgSlice) == 0 {
		return
	}
	command := strings.ToLower(msgSlice[0])

	succeed := b.handleAdminCommand(msg, command)

	if succeed {
		return
	}

	if b.IsMuted() {
		return
	}

	taskSlice := b.parseMessage(msgSlice)

	taskSlice, keysToDelete := b.resolveTasks(taskSlice)

	if len(keysToDelete) > 0 {
		defer b.deleteKeys(keysToDelete)
	}

	finalStreamer := b.assembleStreamer(taskSlice)
	if finalStreamer == nil {
		return
	}

	if b.IsQueueEnabled() {
		b.pushToQueue(finalStreamer)
	} else {
		b.player.Play(finalStreamer)
	}
}

func (b *Bot) handleAdminCommand(msg twitch.PrivateMessage, command string) bool {
	if msg.User.IsBroadcaster || msg.User.IsMod {
		switch command {
		case "mute":
			b.player.Stop()
			b.SetMuted(true)
			b.mtx.Lock()
			b.queue = b.queue[:0]
			b.queueIsPlaying = false
			b.isPlayingSound = false
			b.mtx.Unlock()
			b.printState()
			return true
		case "unmute":
			b.SetMuted(false)
			b.printState()
			return true
		case "qon":
			b.SetQEnabled(true)
			b.printState()
			return true
		case "qoff":
			b.SetQEnabled(false)
			b.printState()
			return true
		case "stop":
			b.player.Stop()
			b.mtx.Lock()
			b.queue = b.queue[:0]
			b.queueIsPlaying = false
			b.isPlayingSound = false
			b.mtx.Unlock()
			return true
		case "skip":
			if b.IsQueueEnabled() {
				b.mtx.Lock()
				if b.isPlayingSound == false {
					b.mtx.Unlock()
					return true
				}
				b.isPlayingSound = false
				b.mtx.Unlock()
			}
			b.player.Stop()
			if b.IsQueueEnabled() {
				b.mtx.Lock()
				if len(b.queue) > 0 {
					b.mtx.Unlock()
					b.playNext()
				} else {
					b.queueIsPlaying = false
					b.mtx.Unlock()
				}
			}
			return true
		case "eron":
			b.mtx.Lock()
			b.erIsOn = true
			b.mtx.Unlock()
			b.printState()
			return true
		case "eroff":
			b.mtx.Lock()
			b.erIsOn = false
			b.mtx.Unlock()
			b.printState()
			return true
		case "vol":
			msgSlice := strings.Fields(msg.Message)
			if len(msgSlice) < 3 {
				return false
			}
			vRaw := msgSlice[2]
			v, err := strconv.Atoi(vRaw)
			if err != nil {
				return false
			}
			if v > 200 {
				v = 200
			}
			if v < 0 {
				v = 0
			}
			b.mtx.Lock()
			b.volume = v
			b.mtx.Unlock()
			b.fileMtx.Lock()
			defer b.fileMtx.Unlock()
			envMap, err := godotenv.Read("config.env")
			if err != nil {
				log.Println("failed to save volume to config.env:", err)
				b.printState()
				return true
			}
			envMap["VOLUME"] = strconv.Itoa(v)
			err = godotenv.Write(envMap, "config.env")
			if err != nil {
				log.Println("failed to save volume to config.env:", err)
			}
			b.printState()
			return true
		}
	}
	return false
}
