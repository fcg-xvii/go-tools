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

func (s ActionData) isResponse() bool {
	_, check := s["Response"]
	return check
}

func actionDataFromRaw(src []byte) (res ActionData) {
	res, lines := make(ActionData), bytes.Split(src[:len(src)], []byte("\r\n"))
	/// todo...
	for _, line := range lines {
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) == 2 {
			res[string(bytes.TrimSpace(parts[0]))] = string(bytes.TrimSpace(parts[1]))
		}
	}
	return
}

func actionsFromRaw(src []byte, acceptCallback func(ActionData)) (res []byte) {
	if bytes.Index(src, []byte("\r\n\r\n")) < 0 {
		return src
	}
	actionsRaw := bytes.Split(src, []byte("\r\n\r\n"))
	log.Println(len(actionsRaw))
	for i := 0; i < len(actionsRaw)-1; i++ {
		action := actionDataFromRaw(actionsRaw[i])
		acceptCallback(action)
	}
	res = actionsRaw[len(actionsRaw)-1]
	return
}

// Request
func initRequest(action string, data ActionData, chanResponse chan Response) Request {
	data["Action"] = action
	return Request{data, chanResponse}
}

type Request struct {
	ActionData
	chanResponse chan Response
}

// Response
func initResponseError(err error) Response {
	return Response{
		ActionData{
			"Action":  "Error",
			"Message": err.Error(),
		},
	}
}

type Response struct {
	ActionData
}

func (s Response) IsError() bool {
	return s.ActionData["Action"] == "Error"
}

func (s Response) ErrorMessage() string {
	return s.ActionData["Message"]
}

// Event
type Event struct {
	ActionData
}
