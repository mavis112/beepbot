package bot

import (
	"testing"

	"github.com/gopxl/beep/v2"
)

func TestConvertMsg(t *testing.T) {
	buf := map[string]*beep.Buffer{
		"car":   nil,
		"alert": nil,
		"a":     nil,
	}
	ttsLang := map[string]string{
		"ru": "ru-RU",
		"en": "en-US",
		"jp": "ja-JP",
	}
	testCases := []struct {
		name  string
		input string
		want  string
		ok    bool
	}{
		{
			name:  "not a command",
			input: "!крыса",
			want:  "",
			ok:    false,
		},
		{
			name:  "2 tts commands",
			input: "!en Hello+!ru Дай мне а",
			want:  "!m en Hello ru- Дай мне а",
			ok:    true,
		},
		{
			name:  "invalid effect",
			input: "!en-something Hello",
			want:  "!m en-something Hello",
			ok:    true,
		},
		{
			name:  "complex cutting with reverse",
			input: "!car-c80-sk50-r",
			want:  "!m car-rs-cs40-ce20",
			ok:    true,
		},
		{
			name:  "complex cutting chain",
			input: "!alert-sk50-c80-sk20-c50",
			want:  "!m alert-cs58-ce26",
			ok:    true,
		},
		{
			name:  "chaos #1",
			input: "!ALARM-ER-f40+!EN-s20 Hello+!rand-C50",
			want:  "!m alarm-er-sp140 en-sp80 Hello rand-ce50",
			ok:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := convertMsg(tc.input, buf, ttsLang)
			if got != tc.want || ok != tc.ok {
				t.Errorf("msg: %s\ngot: %s\n%t\nwant: %s\n%t", tc.input, got, ok, tc.want, tc.ok)
			}
		})
	}
}
