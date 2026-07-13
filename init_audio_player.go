package main

import (
	"beepbot/audio"
	"log"
	"strings"
)

func initAudioPlayer(deviceName string) (audio.AudioOutput, error) {
	if strings.ToLower(deviceName) == "oto" {
		log.Println("Enforcing oto player initialization")
		player, err := audio.NewOtoPlayer(44100)
		if err != nil {
			return nil, err
		}
		log.Println("Successfully initialized oto player")
		return player, nil
	}
	player, err := audio.NewMalgoPlayer(44100, deviceName)
	if err != nil {
		log.Printf("Malgo player failed to start: %v. Falling back to Oto.", err)
		player, err := audio.NewOtoPlayer(44100)
		if err != nil {
			return nil, err
		}
		log.Println("Successfully initialized legacy oto player")
		return player, nil
	}
	log.Println("Successfully initialized malgo player")
	return player, nil
}
