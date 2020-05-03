package yask

import (
	"fmt"
	"io"
	"net/url"
	"strconv"
)

const (
	YaSTTUrl   = "https://stt.api.cloud.yandex.net/speech/v1/stt:recognize"
	FormatLPCM = "lpcm"
	FormatOgg  = "oggopus"
	Rate8k     = "8000"
	Rate16k    = "16000"
	Rate48k    = "48000"
)

type STTConfig struct {
	Lang            string
	Topic           string
	ProfanityFilter bool
	Format          string
	Rate            string
	YaFolderID      string
	YaAPIKey        string
	Data            io.Reader
}

func (s *STTConfig) URI() string {
	vals := url.Values{
		"lang":            s.Lang,
		"topic":           s.Topic,
		"profanityFilter": strconv.FormatBool(s.ProfanityFilter),
		"format":          s.Format,
		"simpleRateHertz": s.Rate,
		"folderId":        s.YaFolderID,
	}

	url := fmt.Sprintf("%v?%v", YaSTTUrl, vals.Encode())
	return url
}

func SpeechToTextShort() {

}
