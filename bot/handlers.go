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
			b.mtx.Lock()
			b.speakerIsMuted = true
			b.queue = b.queue[:0]
			b.queueIsPlaying = false
			b.isPlayingSound = false
			b.playbackID++
			b.mtx.Unlock()
			b.PrintState()
			return true
		case "unmute":
			b.mtx.Lock()
			b.speakerIsMuted = false
			b.mtx.Unlock()
			b.PrintState()
			return true
		case "qon":
			b.mtx.Lock()
			b.queueEnabled = true
			b.mtx.Unlock()
			if err := b.saveEnvParam("QUEUE", "true"); err != nil {
				log.Println("failed to save queue state to config.env:", err)
				b.PrintState()
				return true
			}
			b.PrintState()
			return true
		case "qoff":
			b.mtx.Lock()
			b.queueEnabled = false
			b.mtx.Unlock()
			if err := b.saveEnvParam("QUEUE", "false"); err != nil {
				log.Println("failed to save queue state to config.env:", err)
				b.PrintState()
				return true
			}
			b.PrintState()
			return true
		case "stop":
			b.player.Stop()
			b.mtx.Lock()
			b.queue = b.queue[:0]
			b.queueIsPlaying = false
			b.isPlayingSound = false
			b.playbackID++
			b.mtx.Unlock()
			return true
		case "skip":
			b.mtx.Lock()
			isQEnabled := b.queueEnabled
			if isQEnabled {
				if !b.isPlayingSound {
					b.mtx.Unlock()
					return true
				}
				b.isPlayingSound = false
				b.playbackID++
			}
			b.mtx.Unlock()

			b.player.Stop()
			if isQEnabled {
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
			if err := b.saveEnvParam("ER", "true"); err != nil {
				log.Println("failed to save er effect state to config.env:", err)
				b.PrintState()
				return true
			}
			b.PrintState()
			return true
		case "eroff":
			b.mtx.Lock()
			b.erIsOn = false
			b.mtx.Unlock()
			if err := b.saveEnvParam("ER", "false"); err != nil {
				log.Println("failed to save er effect state to config.env:", err)
				b.PrintState()
				return true
			}
			b.PrintState()
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
			vStr := strconv.Itoa(v)
			if err = b.saveEnvParam("VOLUME", vStr); err != nil {
				log.Println("failed to save volume to config.env:", err)
				b.PrintState()
				return true
			}
			b.PrintState()
			return true
		}
	}
	return false
}

func (b *Bot) saveEnvParam(key, val string) error {
	b.fileMtx.Lock()
	defer b.fileMtx.Unlock()
	envMap, err := godotenv.Read("config.env")
	if err != nil {
		return err
	}
	envMap[key] = val
	if err = godotenv.Write(envMap, "config.env"); err != nil {
		return err
	}
	return nil
}
