package ami

import (
	"bytes"
	"fmt"
	"log"
)

type ActionData map[string]string

func (s ActionData) raw() (res []byte) {
	for key, val := range s {
		res = append(res, []byte(fmt.Sprintf("%v: %v\r\n", key, val))...)
	}
	res = append(res, []byte("\r\n")...)
	return
}

func actionDataFromRaw(src []byte) (res ActionData) {
	res, lines := make(ActionData), bytes.Split(src[:len(src)-4], []byte("\r\n"))
	for _, line := range lines {
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) == 2 {
			res[string(bytes.TrimSpace(parts[0]))] = string(bytes.TrimSpace(parts[1]))
		}
	}
	return
}

func actionsFromRaw(src []byte, chanResponse chan Response, chanEvent chan Event) (res []byte) {
	if bytes.Index(src, []byte("\r\n\r\n")) < 0 {
		return src
	}
	actionsRaw := bytes.Split(src, []byte("\r\n\\r\n"))
	for i := 0; i < len(actionsRaw)-1; i++ {
		action := actionDataFromRaw(actionsRaw[i])
		log.Println(action)
	}
	res = actionsRaw[len(actionsRaw)-1]
	log.Println("RES: ", string(res))
	return
}

// Request
type Request struct {
}

// Response
type Response struct {
}

// Event
type Event struct {
}
