package bufio

import (
	"bytes"
	"log"
	"testing"
)

var (
	dataBuf bytes.Buffer
)

func init() {
	dataBuf.Write([]byte("{ one }\r\n{ two }\r\n{ three }"))
}

func TestDelim(t *testing.T) {
	delim := []byte("}\r\n")
	r := NewReader(&dataBuf)
	for {
		data, err := r.ReadBytes(delim)
		log.Println(string(data), err, string(DelimRemove(data, delim)))
		if err != nil {
			break
		}
	}
}
