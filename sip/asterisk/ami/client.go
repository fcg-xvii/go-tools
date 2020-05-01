package ami

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
	"time"
)

type State byte

const (
	StateStopped State = iota
	StateConnection
	StateConnected
	StateAuth
	StateAvailable
	StateBusy
)

func (s State) String() string {
	switch s {
	case StateStopped:
		return "Stopped"
	case StateConnection:
		return "Connection"
	case StateConnected:
		return "Connected"
	case StateAuth:
		return "Auth"
	case StateAvailable:
		return "Available"
	case StateBusy:
		return "Busy"
	default:
		return ""
	}
}

func New(host, login, password string, stateChanged func(State, error)) (cl *Client) {
	cl = &Client{
		&client{
			host:         host,
			login:        login,
			password:     password,
			stateChanged: stateChanged,
			state:        StateStopped,
			request:      make(chan Request),
			response:     make(chan Response),
			event:        make(chan Event),
			queue:        list.New(),
			socketClosed: make(chan error),
		},
	}
	runtime.SetFinalizer(cl, destroyClient)
	return
}

// Client object
type Client struct {
	*client
}

type client struct {
	host            string
	login           string
	password        string
	conn            net.Conn
	request         chan Request
	response        chan Response
	event           chan Event
	clientSideEvent chan Event
	stateChanged    func(State, error)
	state           State
	queue           *list.List
	socketClosed    chan error
}

//func (s *client)

func (s *client) setState(state State, err error) {
	oldState := s.state
	s.state = state
	if s.stateChanged != nil && (state != oldState || err != nil) {
		s.stateChanged(state, err)
	}
	s.state = state
}

func (s *client) OpenEventChannel() chan Event {
	if s.clientSideEvent == nil {
		s.clientSideEvent = make(chan Event)
	}
	return s.clientSideEvent
}

func (s *client) eventAccepted(event Event) {
	switch event.Name() {
	case "FullyBooted":
		{
			if s.state == StateAuth {
				if s.queue.Len() == 0 {
					s.setState(StateAvailable, nil)
				} else {
					s.sendQueueRequest()
				}
			}
		}
	}

	// send event to client side
	if s.clientSideEvent != nil {
		s.clientSideEvent <- event
	}
}

// start open connection to asterisk ami server
func (s *client) Start() {

	var err error

	// check state. StateStopped needed
	if s.state != StateStopped {
		err = errors.New("AMI start error: client already started")
		s.stateChanged(s.state, err)
		return
	}

	defer func() {
		s.setState(StateStopped, err)
	}()

	s.setState(StateConnection, nil)

	// connection and read ami greetings message
	if s.conn, err = net.Dial("tcp", s.host); err != nil {
		err = fmt.Errorf("AMI connection socket connection error: %v", err.Error())
		return
	}
	s.setState(StateConnected, nil)

	// socket connected. receive greetings text
	if _, err = s.receiveSingle(); err != nil {
		err = fmt.Errorf("AMI greetings receive error: %v", err.Error())
		return
	}

	// greetings received, make attempt to auth
	auth := InitRequest("Login")
	auth.SetParam("UserName", s.login)
	auth.SetParam("Secret", s.password)

	/*ActionData{
			"UserName": s.login,
			"Secret":   s.password,
		},
		make(chan Response),
	)*/

	actionCallback := func(action ActionData) {
		log.Println(action)
		if action.isEvent() {
			s.eventAccepted(Event{action})
		} else {
			response := Response{action}
			if !response.IsError() {
				s.setState(StateAuth, nil)
			} else {
				err = fmt.Errorf("AMI authentication error: %v", action["Message"])
				return
			}
		}
	}

	if socketErr := s.sendSingleRequest(auth, actionCallback); socketErr != nil || err != nil {
		if err == nil {
			err = socketErr
		}
		return
	}

	go s.receiveLoop()

loop:
	for {
		select {
		case request := <-s.request:
			{
				s.queue.PushBack(request)
				if s.state == StateAvailable {
					s.setState(StateBusy, nil)
					if err = s.sendQueueRequest(); err != nil {
						break loop
					}
				}

			}
		case event := <-s.event:
			s.eventAccepted(event)
		case response := <-s.response:
			{
				reqElem := s.queue.Front()
				request := reqElem.Value.(Request)
				request.chanResponse <- response
				close(request.chanResponse)
				s.queue.Remove(reqElem)
				if err = s.sendQueueRequest(); err != nil {
					break loop
				}
			}
		case err = <-s.socketClosed:
			break loop
		}
	}

	return
}

func (s *client) sendQueueRequest() error {
	if s.queue.Len() > 0 {
		s.setState(StateBusy, nil)
		req := s.queue.Front().Value.(Request)
		return s.sendRequest(req)
	} else {
		s.setState(StateAvailable, nil)
		return nil
	}
}

func (s *client) receiveSingle() (data []byte, err error) {
	count, buf := 0, make([]byte, 1024)
	if count, err = s.conn.Read(buf); err == nil {
		data = buf[:count]
	}
	return
}

func (s *client) sendSingleRequest(request Request, acceptCallback func(ActionData)) (err error) {
	// send action
	if err = s.sendRequest(request); err != nil {
		return
	}

	// receive answer
	var data []byte
	for {
		count, buf := 0, make([]byte, 1024)
		if count, err = s.conn.Read(buf); err != nil {
			return
		}
		if data = actionsFromRaw(append(data, buf[:count]...), acceptCallback); len(data) == 0 {
			return
		}
	}
}

func (s *client) sendRequest(request Request) (err error) {
	if _, err := s.conn.Write(request.raw()); err != nil {
		err = fmt.Errorf("AMI socket send data error: %v", err.Error())

	}
	return
}

func (s *client) receiveLoop() {
	var (
		data  []byte
		count int
		err   error
	)
	buf := make([]byte, 1024)
	for {
		if count, err = s.conn.Read(buf); err != nil {
			err = fmt.Errorf("AMI socket receive data error: %v", err.Error())
			s.socketClosed <- err
			return
		}
		data = actionsFromRaw(
			append(data, buf[:count]...),
			func(action ActionData) {
				if action.isEvent() {
					s.event <- Event{action}
				} else {
					s.response <- Response{action}
				}
			},
		)
	}
}

func (s *client) Request(req Request, timeout time.Duration) (resp Response, accepted bool) {
	req.chanResponse = make(chan Response)
	s.request <- req
	if timeout == 0 {
		resp, accepted = <-req.chanResponse
	} else {
		select {
		case resp, accepted = <-req.chanResponse:
			accepted = true
		case <-time.After(timeout):
			break
		}
	}
	return
}

/*func (s *client) receiveResponse() (res Action, err error) {
	var src []byte
	if src, err = s.receive(); err != nil {
		return
	}
	res, err = actionFromRaw(src)
	return
}

func (s client) sendAction(action Action) (response Action, err error) {
	if _, err = s.conn.Write(action.raw()); err != nil {
		return
	}
	response, err = s.receiveResponse()
	return
}

func (s *client) receiveLoop() {
	var data []byte
	count, buf := 0, make([]byte, 1024)
	for {
		if count, err = s.conn.Read(buf); err != nil {
			break
		}
		data = append(data, buf[:count]...)

	}
	s.conn = nil
}

// exec start main goroutine for exec request to asterisk ami
func (s *client) exec() {
	go func() {
		for {
			select {
			case <-s.done:
				{
					s.disconnect()
					return
				}
			case req := <-s.request:
				{
					action := req.(Action)
					if s.conn == nil {
						if err := s.start(); err != nil {
							s.request <- err
						} else if response, err := s.sendAction(action); err != nil {
							s.request <- err
						} else {
							s.request <- response
						}
					}
				}
			}
		}
	}()
}

func (s *client) Request(action string, data ActionData, chanResponse chan Response) {

}*/

// Close finish work with client
func (s *client) Close() {
	if s.state > StateStopped {
		s.conn.Close()
	}
}

// destructor for finalizer
func destroyClient(cl *Client) {
	cl.Close()
}
