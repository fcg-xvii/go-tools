package ami

import (
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
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

func Open(host, login, password string, stateChanged func(State, error)) (cl *Client) {
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
	host         string
	login        string
	password     string
	conn         net.Conn
	request      chan Request
	response     chan Response
	event        chan Event
	stateChanged func(State, error)
	state        State
}

func (s *client) setState(state State, err error) {
	s.state = state
	if s.stateChanged != nil {
		s.stateChanged(state, err)
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

	// start receive data loop
	/*go s.receiveLoop(func(socketError error) {
		err = socketError
	})*/

	auth := initRequest(
		"Login",
		ActionData{
			"UserName": s.login,
			"Secret":   s.password,
		},
		make(chan Response),
	)

	actionCallback := func(action ActionData) {
		if action.isResponse() {

			if action["Response"] == "Success" {
				s.setState(StateAuth, nil)
			} else {
				s.conn.Close()
				err = fmt.Errorf("AMI authentication error: %v", action["Message"])
				return
			}
		} else {
			if action["Event"] == "FullyBooted" {
				if s.state == StateAuth {
					s.setState(StateAvailable, nil)
				}
			}
		}
	}

	if err = s.sendSingleRequest(auth, actionCallback); err != nil {
		err = fmt.Errorf("AMI greetings receive error: %v", err.Error())
		return
	}

	// queue needed...
	/*loop:
	for {
		select {
		case request := <-s.request:
			{

			}
		case event := <-s.event:
			{

			}
		}
	}*/
	return
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

func (s *client) receiveLoop(errCallback func(error)) {
	var (
		data  []byte
		count int
		err   error
	)
	buf := make([]byte, 1024)
	for {
		if count, err = s.conn.Read(buf); err != nil {
			err = fmt.Errorf("AMI socket receive data error: %v", err.Error())
			errCallback(err)
			return
		}
		data = actionsFromRaw(
			append(data, buf[:count]...),
			func(action ActionData) {
				if action.isResponse() {
					log.Println("RESPONSE...")
				} else {
					log.Println("EVENT ACCEPTED")
				}
			},
		)
	}
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
