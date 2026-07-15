package main

import (
	"beepbot/audio"
	"beepbot/bot"
	"beepbot/tts"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/gopxl/beep/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		finalErr := fmt.Errorf("failed to load env file: %w", err)
		exitWithError(finalErr)
	}

	channel := os.Getenv("CHANNEL")
	if channel == "" {
		finalErr := fmt.Errorf("missing required variable 'CHANNEL'")
		exitWithError(finalErr)
	}

	volume, err := strconv.Atoi(os.Getenv("VOLUME"))
	if err != nil {
		volume = 100
	}

	if volume > 200 {
		volume = 200
	}

	if volume < 0 {
		volume = 0
	}

	var wg sync.WaitGroup

	var (
		soundsBuffer map[string]*beep.Buffer
		errors       []error
		errBuff      error
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		soundsBuffer, errors, errBuff = audio.CreateSoundsBuffer()
	}()

	deviceName := os.Getenv("AUDIO_DEVICE")
	var (
		player audio.AudioOutput
		errPl  error
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		player, errPl = initAudioPlayer(deviceName)
	}()

	wg.Wait()

	if errBuff != nil {
		log.Println("sounds folder missing/empty; TTS only active")
	}
	if len(errors) > 0 {
		for _, e := range errors {
			log.Println(e)
		}
	}
	if errPl != nil {
		exitWithError(errPl)
	}
	defer player.Close()

	ttsLanguages := tts.NewTtsLanguages()

	var (
		queue bool = false
		er    bool = true
	)

	q, err := strconv.ParseBool(os.Getenv("QUEUE"))

	if err == nil {
		queue = q
	}

	e, err := strconv.ParseBool(os.Getenv("ER"))

	if err == nil {
		er = e
	}

	b := bot.New(channel, player, soundsBuffer, queue, er, ttsLanguages, volume)
	msgChan := make(chan twitch.PrivateMessage, 500)

	for range 20 {
		go b.HandleLoop(msgChan)
	}

	b.Client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		select {
		case msgChan <- msg:
		default:
		}
	})

	b.Client.OnSelfJoinMessage(func(msg twitch.UserJoinMessage) {
		log.Printf("Successfully joined channel: %s\n", msg.Channel)
		b.PrintState()
	})

	b.Client.Join(b.Channel)
	if err := b.Client.Connect(); err != nil {
		finalErr := fmt.Errorf("irc connection failed: %w", err)
		exitWithError(finalErr)
	}
}
