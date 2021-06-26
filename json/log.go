package json

import (
	"encoding/json"
	"log"
)

func Log(v interface{}) {
	str, _ := json.MarshalIndent(v, "*", "\t")
	log.Println(string(str))
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
