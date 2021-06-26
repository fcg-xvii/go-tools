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

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
