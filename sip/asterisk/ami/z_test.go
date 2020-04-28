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

	cl := Open(host, login, password, nil)

	req := initRequest(
		"Originate",
		ActionData{
			"Channel": "SIP/777",
			"Async":   "yes",
		},
		make(chan Response),
	)

	cl.start(req)

	resp := <-req.chanResponse

	log.Println("RESP", resp)

	/*resp, err := cl.Request(Action{
		"Action":  "Originate",
		"Channel": "sip/777",
		"Data":    "1234567",
		"Async":   "yes",
	})
	log.Println(resp, err)*/
	time.Sleep(time.Second * 30)
	cl.Close()
}
