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

		if !okSound && !okTts {
			return "", false
		}
		convertedMsg := "!m " + convertCommand(message)
		return convertedMsg, true
	}
	messageSlice := strings.Split(message, "+!")
	convertedSlice := make([]string, 0, len(messageSlice))
	hadTts := false
	for _, command := range messageSlice {
		cmdSlice := strings.SplitN(command, " ", 2)
		cmdNameSlice := strings.Split(cmdSlice[0], "-")
		convertedCommand := convertCommand(command)

		if hadTts && len(cmdNameSlice) == 1 {
			convertedCommand += "-"
		}

		convertedSlice = append(convertedSlice, convertedCommand)

		_, okSound := buf[strings.ToLower(cmdNameSlice[0])]
		_, okTts := ttsLang[strings.ToLower(cmdNameSlice[0])]
		if okTts {
			hadTts = true
		}
		if okSound {
			hadTts = false
		}
	}
	convertedMsg := "!m " + strings.Join(convertedSlice, " ")
	return convertedMsg, true
}

func convertCommand(msg string) string {
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
		convertedEffects = append(convertedEffects, "cs"+strconv.Itoa(int(startPercent)))
		convertedEffects = append(convertedEffects, "ce"+strconv.Itoa(int(100.0-endPercent)))

	}
	convertedCmd := name
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
	case "r", "rs":
		return "rs"
	case "rv", "dl":
		return "dl"
	case "er", "lq", "vb", "ga", "tr", "ts":
		return loweredEffName
	case "rm", "st", "sp", "cs", "ce":
		return loweredEffName + effValue
	case "f":
		if effValue == "" {
			return "sp130"
		}

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
		if effValue == "" {
			return "sp70"
		}
		v, err := strconv.Atoi(effValue)
		if err != nil {
			return "sp70"
		}
		v = 100 - v
		if v < 10 {
			v = 10
		}
		return "sp" + strconv.Itoa(v)
	default:
		return ""
	}
}

func parseEff(eff string) (effName, effValue string) {
	digitId := -1
	for i, r := range eff {
		if r >= '0' && r <= '9' {
			digitId = i
			break
		}
	}
	effName = eff
	if digitId >= 0 {
		effName = eff[:digitId]
		effValue = eff[digitId:]
	}
	return effName, effValue
}
