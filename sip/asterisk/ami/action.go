package ami

import (
	"bytes"
	"errors"
	"fmt"
)

type Action map[string]string

func (s Action) raw() (res []byte) {
	for key, val := range s {
		res = append(res, []byte(fmt.Sprintf("%v: %v\r\n", key, val))...)
	}
	res = append(res, []byte("\r\n")...)
	return
}

func actionFromRaw(src []byte) (res Action, err error) {
	if !bytes.HasSuffix(src, []byte("\r\n\r\n")) {
		err = errors.New("Expected action suffix '\r\n\r\n'")
	} else {
		res = make(Action)
		lines := bytes.Split(src[:len(src)-4], []byte("\r\n"))
		for _, line := range lines {
			parts := bytes.SplitN(line, []byte(":"), 2)
			if len(parts) == 2 {
				res[string(bytes.TrimSpace(parts[0]))] = string(bytes.TrimSpace(parts[1]))
			}
		}
	}
	return
}

func actionsFromRaw(src []byte, event, action chan Action) []byte {
	actionsRaw := bytes.Split(src, "\r\n")
}
