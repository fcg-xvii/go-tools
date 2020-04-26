package ami

import (
	"log"
	"net"
	"runtime"
)

func Open(host, login, password string, eventListener chan Action) (cl *Client) {
	cl = &Client{
		&client{
			host:           host,
			login:          login,
			password:       password,
			done:           make(chan struct{}),
			request:        make(chan Action),
			response:       make(chan interface{}),
			event:          make(chan Action),
			eventOtherSide: eventListener,
		},
	}
	runtime.SetFinalizer(cl, destroyClient)
	cl.exec()
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
	request        chan Action
	response       chan interface{}
	event          chan Action
	eventOtherSide chan Action
	started        bool
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
		return
	}
	if err = s.receiveGreetings(); err != nil {
		return
	}

	// connection accepted, send auth data
	response, err := s.sendAction(Action{
		"Action":   "Login",
		"Username": s.login,
		"Secret":   s.password,
	})
	log.Println(response, err)
	return
}

func (s *client) receiveGreetings() (err error) {
	_, err = s.receive()
	return
}

func (s *client) receiveResponse() (res Action, err error) {
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

func (s *client) receive() {
	var data []byte
	count, buf := 0, make([]byte, 1024)
	for {
		if count, err = s.conn.Read(buf); err != nil {
			break
		}
		data = append(data, buf[:count]...)

	}
	s.conn = nil
	/*count, buf := 0, make([]byte, 1024)
	for {
		s.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 50))
		if count, err = s.conn.Read(buf); err != nil {
			if e, ok := err.(interface{ Timeout() bool }); ok && e.Timeout() {
				log.Println("RECEIVED", string(res))
				err = nil
			} else {
				s.conn = nil
			}
			return
		} else {
			res = append(res, buf[:count]...)
		}
	}*/
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

func (s *client) Request(action Action) (response Action, err error) {
	s.request <- action
	resp := <-s.request
	if obj, isErr := resp.(error); isErr {
		err = obj
	} else {
		response = resp.(Action)
	}
	return
}

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
