package json

import (
	"encoding/json"
	"io"
	"log"
	"os"
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

func UnmarshalReader(r io.Reader, v interface{}) error {
	return Decode(r, v)
}

func UnmarshalFile(fileName string, v interface{}) (err error) {
	var f *os.File
	if f, err = os.Open(fileName); err != nil {
		return
	}
	err = Decode(f, v)
	f.Close()
	return
}
