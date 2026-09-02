package bot

import (
	"strconv"
	"strings"

	"github.com/gopxl/beep/v2"
)

func convertMsg(msg string, buf map[string]*beep.Buffer, ttsLang map[string]string) (string, bool) {
	message, _ := strings.CutPrefix(msg, "!")

	if !strings.Contains(message, "+!") {
		command := strings.Split(message, " ")[0]
		commandName := strings.ToLower(strings.Split(command, "-")[0])
		_, okSound := buf[commandName]
		_, okTts := ttsLang[commandName]

		if !okSound && !okTts && commandName != "rand" {
			return "", false
		}
		convertedMsg := "!m " + convertCommand(message, false)
		return convertedMsg, true
	}
	messageSlice := strings.Split(message, "+!")
	convertedSlice := make([]string, 0, len(messageSlice))
	hadTts := false
	for _, command := range messageSlice {
		cmdSlice := strings.SplitN(command, " ", 2)
		cmdNameSlice := strings.Split(cmdSlice[0], "-")
		convertedCommand := convertCommand(command, hadTts)

		convertedSlice = append(convertedSlice, convertedCommand)
		commandName := strings.ToLower(cmdNameSlice[0])
		_, okSound := buf[commandName]
		_, okTts := ttsLang[commandName]
		if okTts {
			hadTts = true
		}
		if okSound || commandName == "rand" {
			hadTts = false
		}
	}
	convertedMsg := "!m " + strings.Join(convertedSlice, " ")
	return convertedMsg, true
}

func convertCommand(msg string, hadTts bool) string {
	msgSlice := strings.SplitN(msg, " ", 2)
	var text string
	if len(msgSlice) > 1 {
		text = msgSlice[1]
	}
	cmd := msgSlice[0]
	cmdSlice := strings.Split(cmd, "-")
	name := strings.ToLower(cmdSlice[0])
	convertedEffects := make([]string, 0, len(cmdSlice))

	startPercent := 0.0
	endPercent := 100.0
	var (
		currentLen  float64
		hasClipping bool
	)
	for _, eff := range cmdSlice[1:] {
		effName, effValue := parseEff(eff)
		switch strings.ToLower(effName) {
		case "c":
			var (
				v   int
				err error
			)
			v, err = strconv.Atoi(effValue)
			if err != nil {
				v = 50
			}
			val := float64(v)
			currentLen = endPercent - startPercent
			endPercent = startPercent + currentLen*(val/100.0)
			hasClipping = true
			continue
		case "sk":
			var (
				v   int
				err error
			)
			v, err = strconv.Atoi(effValue)
			if err != nil {
				v = 50
			}
			val := float64(v)
			currentLen = endPercent - startPercent
			startPercent = startPercent + currentLen*(val/100.0)
			hasClipping = true
			continue

		}
		convertedEffect := convertEffects(eff)
		if convertedEffect != "" {
			convertedEffects = append(convertedEffects, convertedEffect)
		}
	}
	if hasClipping {
		if startPercent > 0 {
			convertedEffects = append(convertedEffects, "cs"+strconv.Itoa(int(startPercent)))
		}
		if endPercent < 100 {
			convertedEffects = append(convertedEffects, "ce"+strconv.Itoa(int(100.0-endPercent)))
		}
	}
	convertedCmd := name
	if hadTts && len(convertedEffects) == 0 {
		convertedCmd += "-"
	}
	if len(convertedEffects) > 0 {
		convertedCmd += "-" + strings.Join(convertedEffects, "-")
	}

	if text != "" {
		convertedCmd += " " + text
	}
	return convertedCmd
}

func convertEffects(eff string) string {
	effName, effValue := parseEff(eff)
	loweredEffName := strings.ToLower(effName)
	switch loweredEffName {
	case "r":
		return "rs"
	case "f":
		v, err := strconv.Atoi(effValue)
		if err != nil {
			return "sp130"
		}
		v += 100
		if v > 200 {
			v = 200
		}
		return "sp" + strconv.Itoa(v)
	case "s":
		v, err := strconv.Atoi(effValue)
		if err != nil {
			return "sp70"
		}
		v = 100 - v
		if v < 10 {
			v = 10
		}
		return "sp" + strconv.Itoa(v)
	}
	return loweredEffName + effValue
}

func parseEff(eff string) (effName, effValue string) {
	digitId := -1
	for i, r := range eff {
		if r >= '0' && r <= '9' {
			digitId = i
			break
		}
	}
	if digitId >= 0 {
		effName = eff[:digitId]
		effValue = eff[digitId:]
	} else {
		effName = eff
	}
	return effName, effValue
}
