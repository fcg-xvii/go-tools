package ami

import (
	"fmt"

	"github.com/fcg-xvii/go-tools/json/jsonmap"
)

func InitRequest(action string) Request {
	return Request{
		ActionData: ActionData{
			"Action": action,
		},
	}
}

type Request struct {
	ActionData
	Variables    jsonmap.JSONMap
	chanResponse chan Response
	sended       bool
}

func (s *Request) SetParam(key, value string) {
	if len(value) > 0 {
		s.ActionData[key] = value
	}
}

func (s *Request) SetVariable(key, value string) {
	if s.Variables == nil {
		s.Variables = make(jsonmap.JSONMap)
	}
	s.Variables[key] = value
}

func (s *Request) SetVariables(m jsonmap.JSONMap) {
	if s.Variables == nil {
		s.Variables = make(jsonmap.JSONMap)
	}
	for key, val := range m {
		s.Variables[key] = val
	}
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
