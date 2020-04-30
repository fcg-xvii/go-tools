package ami

import (
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
	log.Println(host, login, password)

	var cl *Client
	cl = New(host, login, password, func(state State, err error) {
		log.Println("STATE_CHANGED", state, err)
		switch state {
		case StateStopped:
			{
				time.Sleep(time.Second * 5)
				log.Println("Reconnect...")
				cl.Start()
			}
		}
	})

	go cl.Start()

	req := InitRequest("Originate")
	req.SetParam("Channel", "sip/777")
	req.SetParam("Async", "yes")
	req.SetVariable("one", "1")
	req.SetVariable("two", "2")

	resp, accepted := cl.Request(req, 0)

	log.Println("RESP!!!!!!!!!!!!!!!!!!!!!!!!!", resp, accepted)

	//log.Println(resp, err)

	time.Sleep(time.Second * 10)
	cl.Close()
}
