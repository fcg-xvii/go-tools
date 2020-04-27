package ami

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"time"
)

func Open(host, login, password string, eventListener chan Event) (cl *Client) {
	cl = &Client{
		&client{
			host:           host,
			login:          login,
			password:       password,
			done:           make(chan struct{}),
			request:        make(chan Request),
			response:       make(chan interface{}),
			event:          make(chan Event),
			eventOtherSide: eventListener,
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
	host           string
	login          string
	password       string
	conn           net.Conn
	done           chan struct{}
	request        chan Request
	response       chan interface{}
	event          chan Event
	eventOtherSide chan Event
}

// exist close connection if opened
func (s *client) disconnect() {
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
}

// start open connection to asterisk ami server
func (s *client) start() (err error) {
	// connection and read ami greetings message
	if s.conn, err = net.Dial("tcp", s.host); err != nil {
		err = fmt.Errorf("AMI connection socket connection error: %v", err.Error())
		return
	}

	// socket connected. receive greetings text
	var greetings []byte
	if greetings, err = s.receiveSingle(); err != nil {
		err = fmt.Errorf("AMI greetings receive error: %v", err.Error())
		return
	}

	log.Println("GREETINGS ACCEPTED", string(greetings))
	/*
		// connection accepted, send auth data
		response, err := s.sendAction(Action{
			"Action":   "Login",
			"Username": s.login,
			"Secret":   s.password,
		})
		log.Println(response, err)
	*/
	return
}

func (s *client) receiveSingle() (data []byte, err error) {
	count, buf := 0, make([]byte, 1024)
	for {
		s.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 50))
		if count, err = s.conn.Read(buf); err != nil {
			if e, ok := err.(interface{ Timeout() bool }); ok && e.Timeout() {
				log.Println("RECEIVED", string(data))
				err = nil
			} else {
				s.conn = nil
			}
			return
		} else {
			data = append(data, buf[:count]...)
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
	if s.done != nil {
		close(s.done)
	}
}

// destructor for finalizer
func destroyClient(cl *Client) {
	close(cl.done)
}
