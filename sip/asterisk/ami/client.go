package ami

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"
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

func New(host, login, password string, ctxGlobal context.Context, stateChanged func(State, error)) (cl *Client) {
	cl = &Client{
		&client{
			host:           host,
			login:          login,
			password:       password,
			stateChanged:   stateChanged,
			state:          StateStopped,
			request:        make(chan Request),
			response:       make(chan Response),
			event:          make(chan Event),
			requestsWork:   list.New(),
			socketClosed:   make(chan error),
			actionIDPrefix: fmt.Sprint(time.Now().UnixNano()),
			eventListeners: make(map[int64]*EventListener),
			locker:         new(sync.RWMutex),
		},
	}
	if ctxGlobal == nil {
		cl.ctx, cl.ctxCancel = context.WithCancel(context.Background())
	} else {
		cl.ctx, cl.ctxCancel = context.WithCancel(ctxGlobal)
	}
	go cl.eventListenersCleaner()
	runtime.SetFinalizer(cl, destroyClient)
	return
}

// Client object
type Client struct {
	*client
}

type client struct {
	ctx             context.Context
	ctxCancel       context.CancelFunc
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
	requestsWork    *list.List
	socketClosed    chan error
	actionIDPrefix  string
	actionUUID      uint64
	eventListeners  map[int64]*EventListener
	locker          *sync.RWMutex
}

func (s *client) removeEventListener(uuid int64) {
	s.locker.Lock()
	if listener, check := s.eventListeners[uuid]; check {
		listener.close()
		delete(s.eventListeners, uuid)
	}
	s.locker.Unlock()
}

func (s *client) eventListenersCleaner() {
	ctx, _ := context.WithCancel(s.ctx)
	for {
		select {
		case <-time.After(time.Minute * 30):
			{
				now := time.Now()
				s.locker.Lock()
				for uuid, v := range s.eventListeners {
					if now.After(v.timeActual) {
						v.close()
						delete(s.eventListeners, uuid)
					}
				}
				s.locker.Unlock()
			}
		case <-ctx.Done():
			{
				for _, v := range s.eventListeners {
					v.close()
				}
				return
			}
		}
	}
}

func (s *client) registerEventListener(uuid int64) <-chan Event {
	listener := &EventListener{
		uuid:      uuid,
		eventChan: make(chan Event),
	}
	s.locker.Lock()
	s.eventListeners[uuid] = listener
	s.locker.Unlock()
	return listener.eventChan
}

func (s *client) initActionID() (res string) {
	res = fmt.Sprintf("%v%v", s.actionIDPrefix, s.actionUUID)
	if s.actionUUID < max_client_uuid {
		s.actionUUID++
	} else {
		s.actionUUID = 0
	}
	return
}

func (s *client) requestByActionID(actionID string) (req Request, elem *list.Element, check bool) {
	for elem = s.requestsWork.Front(); elem != nil; elem = elem.Next() {
		req = elem.Value.(Request)
		if req.ActionID() == actionID {
			check = true
			return
		}
	}
	return
}

func (s *client) setState(state State, err error) {
	oldState := s.state
	s.state = state
	if s.stateChanged != nil && (state != oldState || err != nil) {
		s.stateChanged(state, err)
	}
	s.state = state
}

func (s *client) Event() chan Event {
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
				for elem := s.requestsWork.Front(); elem != nil; elem.Next() {
					s.sendQueueRequest()
				}
			}
		}
	}

	// send event to client side
	if s.clientSideEvent != nil {
		s.clientSideEvent <- event
	}

	if event.uuid > 0 {
		var check bool
		var listener *EventListener
		s.locker.RLock()
		if listener, check = s.eventListeners[event.uuid]; check {
			listener.incomingEvent(event)
		}
		s.locker.RUnlock()
		if check && event.Name() == "Hangup" {
			s.locker.Lock()
			listener.close()
			delete(s.eventListeners, event.uuid)
			s.locker.Unlock()
		}
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

	actionCallback := func(action ActionData) {
		if action.isEvent() {
			s.eventAccepted(Event{action, 0})
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
				actionID := s.initActionID()
				request.ActionData["ActionID"] = actionID
				s.requestsWork.PushFront(request)
				if s.state == StateAuth {
					if err := s.sendRequest(request); err != nil {
						log.Println("SendRequestERROR")
					}
				}
			}
		case event := <-s.event:
			s.eventAccepted(event)
		case response := <-s.response:
			{
				if req, elem, check := s.requestByActionID(response.ActionID()); check {
					req.chanResponse <- response
					close(req.chanResponse)
					s.requestsWork.Remove(elem)
				}
			}
		case err = <-s.socketClosed:
			break loop
		}
	}

	return
}

func (s *client) sendQueueRequest() error {
	for elem := s.requestsWork.Front(); elem != nil; elem.Next() {
		req := elem.Value.(Request)
		if req.sended {
			s.requestsWork.Remove(elem)
		} else if err := s.sendRequest(req); err != nil {
			return err
		}
	}
	return nil
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

func (s *client) sendRequest(req Request) (err error) {
	if _, err := s.conn.Write(req.raw()); err != nil {
		err = fmt.Errorf("AMI socket send data error: %v", err.Error())
	} else {
		req.sended = true
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
					s.event <- initEvent(action)
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

// Close finish work with client
func (s *client) Close() {
	if s.state > StateStopped {
		s.conn.Close()
	}
	s.ctxCancel()
}

// destructor for finalizer
func destroyClient(cl *Client) {
	cl.Close()
}
