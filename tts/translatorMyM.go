package tts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type respMyM struct {
	ResponseData translatedText `json:"responseData"`
}

type translatedText struct {
	TranslatedText string `json:"translatedText"`
}

func getTranslateUrlMym(lang, text string) string {
	v := url.Values{}
	v.Add("langpair", "Autodetect|"+lang)
	v.Add("q", text)
	return "https://api.mymemory.translated.net/get?" + v.Encode()
}

func TranslateMym(lang, text string) (string, error) {
	url := getTranslateUrlMym(lang, text)
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status not ok: %d", resp.StatusCode)
	}

	var result respMyM

	limitData := io.LimitReader(resp.Body, 8<<10)

	data, err := io.ReadAll(limitData)
	if err != nil {
		return "", fmt.Errorf("can't read data: %v", err)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("can't decode data: %v", err)
	}

	return result.ResponseData.TranslatedText, nil
}
