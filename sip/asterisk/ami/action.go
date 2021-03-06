package ami

import (
	"bytes"
	"fmt"
)

type ActionData map[string]string

func (s ActionData) raw() (res []byte) {
	for key, val := range s {
		res = append(res, []byte(fmt.Sprintf("%v: %v\r\n", key, val))...)
	}
	res = append(res, []byte("\r\n")...)
	return
}

func (s ActionData) isEvent() bool {
	_, check := s["Event"]
	return check
}

func (s ActionData) ActionID() string {
	return s["ActionID"]
}

func actionDataFromRaw(src []byte) (res ActionData) {
	res, lines := make(ActionData), bytes.Split(src, []byte("\r\n"))
	/// todo...
	for _, line := range lines {
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) == 2 {
			res[string(bytes.TrimSpace(parts[0]))] = string(bytes.TrimSpace(parts[1]))
		}
	}
	return
}

func actionsFromRaw(src []byte, accept func(ActionData)) (res []byte) {
	if bytes.Index(src, []byte("\r\n\r\n")) < 0 {
		return src
	}
	actionsRaw := bytes.Split(src, []byte("\r\n\r\n"))
	for i := 0; i < len(actionsRaw)-1; i++ {
		action := actionDataFromRaw(actionsRaw[i])
		accept(action)
	}
	res = actionsRaw[len(actionsRaw)-1]
	return
}
