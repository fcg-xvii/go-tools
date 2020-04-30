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

// Request
func InitRequest(action string) Request {
	return Request{
		ActionData: ActionData{
			"Action": action,
		},
	}
}

type Request struct {
	ActionData
	Variables    map[string]string
	chanResponse chan Response
}

func (s *Request) SetParam(key, value string) {
	s.ActionData[key] = value
}

func (s *Request) SetVariable(key, value string) {
	if s.Variables == nil {
		s.Variables = make(map[string]string)
	}
	s.Variables[key] = value
}

func (s *Request) raw() []byte {
	if len(s.Variables) > 0 {
		vars, count := "", 0
		for key, val := range s.Variables {
			vars += fmt.Sprintf("%v=%v", key, val)
			if count < len(s.Variables)-1 {
				vars += "," // todo 1.5 or lower splitter is '|'
			}
			count++
		}
		s.ActionData["Variable"] = vars
	}
	return s.ActionData.raw()
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
	return s.ActionData["Response"] == "Error"
}

func (s Response) ErrorMessage() string {
	return s.ActionData["Message"]
}

// Event
type Event struct {
	ActionData
}

func (s Event) Name() string {
	return s.ActionData["Event"]
}
