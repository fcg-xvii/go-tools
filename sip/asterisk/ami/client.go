package ami

import (
	"io/ioutil"
	"log"
	"net"
	"runtime"
)

func Open(host, login, password string) (cl *Client) {
	cl = &Client{
		&client{
			host:     host,
			login:    login,
			password: password,
			done:     make(chan struct{}),
			request:  make(chan Action),
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
	host     string
	login    string
	password string
	conn     net.Conn
	done     chan struct{}
	request  chan Action
	started  bool
}

// exist close connection if opened
func (s *client) exit() {
	if s.conn != nil {
		s.conn.Close()
	}
}

// start open connection to asterisk ami server
func (s *client) start() (err error) {
	log.Println("START")
	if s.conn, err = net.Dial("tcp", s.host); err != nil {
		return
	}
	s.send(Action{
		"Action":   "Login",
		"Username": s.login,
		"Secret":   s.password,
	})

	return
}

func (s *client) send(action Action) {
	log.Println("SEND_ACTION", action, string(action.raw()))
	_, err := s.conn.Write(action.raw())
	if err != nil {
		log.Println(err)
		return
	}
	s.receive()
}

func (s *client) receive() {
	log.Println("RECEIVE")
	src, err := ioutil.ReadAll(s.conn)
	log.Println(src, err)
}

// exec start main goroutine for exec request to asterisk ami
func (s *client) exec() {
	go func() {
		for {
			select {
			case <-s.done:
				{
					s.exit()
					return
				}
			case action := <-s.request:
				{
					if s.conn == nil {
						if err := s.start(); err != nil {
							log.Println(err)
						}
					}
					log.Println(action)
				}

			}
		}
	}()
}

func (s *client) Request(action Action) {
	s.request <- action
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
