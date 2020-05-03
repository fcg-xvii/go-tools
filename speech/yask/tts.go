package yask

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type TTSConfig struct {
	Text       string
	SSML       string
	Lang       string
	Voice      string
	Emotion    string
	Speed      float32
	Format     string
	Rate       int
	YaFolderID string
	YaAPIKey   string
}

func (s *TTSConfig) isSSML() bool {
	return len(s.SSML) > 0
}

func defaultTTSConfig(yaFolderID, yaAPIKey string) *TTSConfig {
	return &TTSConfig{
		Lang:       LangRU,
		Voice:      VoiceOksana,
		Emotion:    EmotionNeutral,
		Speed:      SpeedStandard,
		Format:     FormatLPCM,
		Rate:       Rate8k,
		YaFolderID: yaFolderID,
		YaAPIKey:   yaAPIKey,
	}
}

// TTsDefaultConfigText returns config with default parameters for raw text recognition and use in TextToSpeech method
func TTSDefaultConfigText(yaFolderID, yaAPIKey, text string) *TTSConfig {
	conf := defaultTTSConfig(yaFolderID, yaAPIKey)
	conf.Text = text
	return conf
}

// TTsDefaultConfigSSML returns config with default parameters for raw text recognition and use in TextToSpeech method
// more details of SSML language in https://cloud.yandex.ru/docs/speechkit/tts/ssml
func TTSDefaultConfigSSML(yaFolderID, yaAPIKey, SSML string) *TTSConfig {
	conf := defaultTTSConfig(yaFolderID, yaAPIKey)
	conf.SSML = SSML
	return conf
}

//
func TextToSpeech(config *TTSConfig) (io.ReadCloser, error) {
	httpForm := url.Values{
		"lang":            []string{config.Lang},
		"voice":           []string{config.Voice},
		"emotion":         []string{config.Emotion},
		"speed":           []string{strconv.FormatFloat(float64(config.Speed), 'f', 1, 32)},
		"format":          []string{config.Format},
		"sampleRateHertz": []string{strconv.FormatInt(int64(config.Rate), 10)},
		"folderId":        []string{config.YaFolderID},
	}
	if config.isSSML() {
		httpForm.Set("ssml", config.SSML)
	} else {
		httpForm.Set("text", config.Text)
	}

	request, err := http.NewRequest("POST", YaTTSUrl, strings.NewReader(httpForm.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", fmt.Sprintf("Api-Key %v", config.YaAPIKey))

	client := new(http.Client)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = unmarshallYaError(response.Body)
		response.Body.Close()
		return nil, err
	}

	return response.Body, nil
}

func unmarshallYaError(r io.Reader) (err error) {
	var data []byte
	if data, err = ioutil.ReadAll(r); err != nil {
		return
	}
	mErr := make(map[string]interface{})
	if err = json.Unmarshal(data, &mErr); err == nil {
		err = fmt.Errorf("Yandex request error: %v", mErr["error_message"])
	}
	return
}

func EncodePCMToWav(in io.Reader, out io.WriteSeeker, sampleRate, bitDepth, numChans int) error {
	encoder := wav.NewEncoder(out, sampleRate, bitDepth, numChans, 1)

	audioBuf := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: numChans,
			SampleRate:  sampleRate,
		},
	}

	for {
		var sample int16
		if err := binary.Read(in, binary.LittleEndian, &sample); err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		audioBuf.Data = append(audioBuf.Data, int(sample))
	}

	if err := encoder.Write(audioBuf); err != nil {
		return err
	}

	return encoder.Close()
}
