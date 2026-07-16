package tts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func NeedTranslate(effects string) (string, bool) {
	effSlice := strings.Split(effects, "-")
	for i, eff := range effSlice {
		if strings.ToLower(eff) == "tr" {
			effSlice[i] = effSlice[len(effSlice)-1]
			effSlice = effSlice[:len(effSlice)-1]
			return strings.Join(effSlice, "-"), true
		}
	}
	return effects, false
}

func getTranslateUrl(lang, text string) string {
	v := url.Values{}
	v.Add("client", "gtx")
	v.Add("dt", "t")
	v.Add("sl", "auto")
	v.Add("tl", lang)
	v.Add("q", text)
	return "https://translate.googleapis.com/translate_a/single?" + v.Encode()
}

func getTranslateReq(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Referer", "https://translate.google.com/")

	return req, nil
}

func Translate(lang, text string) (string, error) {
	reqUrl := getTranslateUrl(lang, text)
	req, err := getTranslateReq(reqUrl)
	if err != nil {
		return "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		if resp == nil {
			return
		}
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 8<<10))
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status not ok: %d", resp.StatusCode)
	}
	var result []any

	limitData := io.LimitReader(resp.Body, 8<<10)

	data, err := io.ReadAll(limitData)
	if err != nil {
		return "", fmt.Errorf("can't read data: %v", err)
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return "", fmt.Errorf("can't decode data: %v", err)
	}
	if len(result) == 0 {
		return "", errors.New("empty translate response")
	}
	var finalResult string

	firstLevel, ok := result[0].([]any)
	if !ok {
		return "", errors.New("bad json structure")
	}

	if len(firstLevel) == 0 {
		return "", errors.New("empty translate response")
	}
	for i := range firstLevel {
		secondLevel, ok := firstLevel[i].([]any)
		if !ok {
			continue
		}

		if len(secondLevel) == 0 {
			continue
		}

		res, ok := secondLevel[0].(string)
		if !ok {
			continue
		}

		finalResult += res
	}

	return finalResult, nil
}
