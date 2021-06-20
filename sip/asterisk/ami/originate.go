package ami

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/fcg-xvii/go-tools/json/jsonmap"
)

func (s *Client) Originate(req *OriginateRequest) (*Originate, error) {
	req.uuid = time.Now().UnixNano()
	timeout := RequestTimeoutDefault
	if req.Timeout > timeout {
		timeout = req.Timeout + time.Millisecond*500
	}

	resp, check := s.Request(req.Request(), timeout)
	if !check {
		return nil, errors.New("Originate request timeout")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("Originate error: %v", resp.ErrorMessage())
	}

	res := initOriginate(req, s)

	return res, nil
}

//////////////////////////////////////////////////////////////////

type OriginateRequest struct {
	Channel     string
	Context     string
	Exten       string
	Priority    string
	Timeout     time.Duration
	CallerID    string
	Variable    jsonmap.JSONMap
	Account     string
	Application string
	Data        string
	uuid        int64
}

func (s *OriginateRequest) Request() (res Request) {
	res = InitRequest("Originate")
	res.SetParam("Channel", s.Channel)
	res.SetParam("Context", s.Context)
	res.SetParam("Exten", s.Exten)
	if s.Timeout > 0 {
		res.SetParam("Timeout", fmt.Sprintf("%v", int(s.Timeout/time.Millisecond)))
	}
	res.SetParam("Priority", s.Priority)
	res.SetParam("CallerID", s.CallerID)
	res.SetParam("Account", s.Account)
	res.SetParam("Application", s.Application)
	res.SetParam("Data", s.Data)
	res.SetParam("Async", "true")
	res.SetParam("ChannelID", fmt.Sprint(s.uuid))
	res.SetVariables(s.Variable)
	return res
}

/////////////////////////////////////////////////////////////////

func initOriginate(req *OriginateRequest, client *Client) *Originate {
	res := &Originate{
		OriginateRequest: req,
		eventChan:        client.registerEventListener(req.uuid),
		locker:           new(sync.RWMutex),
		client:           client,
	}
	go res.listenEvents()
	return res
}

type Originate struct {
	*OriginateRequest
	eventChan      <-chan Event
	userEventChan  chan Event
	locker         *sync.RWMutex
	finished       bool
	err            error
	client         *Client
	responseReason byte
	hangupCause    byte
}

func (s *Originate) listenEvents() {
	for {
		e, ok := <-s.eventChan
		if !ok {
			s.finished = true
			close(s.userEventChan)
			return
		}
		s.locker.RLock()
		if s.userEventChan != nil {
			s.userEventChan <- e
		}
		s.locker.RUnlock()
		switch e.Name() {
		case "OriginateResponse":
			if reason, check := e.ActionData["Reason"]; check {
				reasonVal, _ := strconv.ParseInt(reason, 10, 32)
				s.responseReason = byte(reasonVal)
			}
			log.Println("RREASON", s.responseReason)
			//if s.responseReason != 4 {
			//s.client.removeEventListener(s.uuid)
			//}
		case "Hangup":
			if cause, check := e.ActionData["Cause"]; check {
				causeVal, _ := strconv.ParseInt(cause, 10, 32)
				s.hangupCause = byte(causeVal)
			}
		}
	}
}

func (s *Originate) IsFinished() bool {
	return s.finished
}

func (s *Originate) Events() (res <-chan Event) {
	s.locker.Lock()
	if s.userEventChan == nil {
		s.userEventChan = make(chan Event)
	}
	res = s.userEventChan
	s.locker.Unlock()
	return
}
