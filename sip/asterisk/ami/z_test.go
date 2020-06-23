package ami

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/fcg-xvii/go-tools/text/config"
)

var (
	host, login, password string
)

func init() {
	// setup auth vars, z_auth.config file example :
	// 127.0.0.1:5038::admin::mypassword
	config.SplitFileToVals("z_auth.config", "::", &host, &login, &password)
}

func TestClient(t *testing.T) {
	if host == "" {
		return
	}
	log.Println(host, login, password)

	var cl *Client
	cl = New(host, login, password, nil, func(state State, err error) {
		log.Println("STATE_CHANGED", state, err)
		switch state {
		case StateStopped:
			{
				time.Sleep(time.Second * 5)
				log.Println("Reconnect...")
				go cl.Start()
			}
		}
	})

	/*go func() {
		for e := range cl.Event() {
			log.Println("EVENT", e.Name(), "***", e.ActionData["Uniqueid"], "======", e)
		}
	}()*/

	go cl.Start()

	//req := InitRequest("Originate")
	//req.SetParam("Channel", "sip/777")
	//req.SetParam("Context", "from-test")
	//req.SetParam("Async", "yes")
	//req.SetVariable("one", "1")
	//req.SetVariable("two", "2")

	//resp, accepted := cl.Request(req, 0)

	if originate, err := cl.Originate(&OriginateRequest{
		Channel:  "sip/777",
		Priority: "1",
		Exten:    "s",
		Context:  "call-test",
		Async:    true,
		CallerID: "777",
		Timeout:  time.Second * 15,
	}); err == nil {
		ctx, _ := context.WithCancel(originate.Context())
		ech := originate.Events()
		for {
			select {
			case e := <-ech:
				{
					log.Println("CEVENT", e.Name())
				}
			case <-ctx.Done():
				{
					log.Println("CALL FINISHED")
					return
				}
			}
		}
	} else {
		t.Fatal(err)
	}

	//log.Println("ORIGINATE", originate, err)

	/*cl.Originate(&OriginateRequest{
		Channel:  "sip/101",
		Priority: "1",
		Exten:    "s",
		Context:  "call-test",
		Async:    true,
	})*/

	/*cl.Originate(&OriginateRequest{
		Channel:  "sip/777",
		Priority: "1",
		Exten:    "s",
		Context:  "call-test",
		Async:    true,
		CallerID: "Test-Caller",
		ActionID: "Test-Action",
	})*/

	/*cl.Originate(&OriginateRequest{
		Channel:  "sip/101",
		Context:  "call-test",
		Async:    true,
		CallerID: "Test-Caller",
		ActionID: "Test-Action",
	})*/

	//log.Println("RESP!!!!!!!!!!!!!!!!!!!!!!!!!", resp, accepted)

	//log.Println(resp, err)

	//time.Sleep(time.Second * 300)
	//cl.Close()
}
