package ami

import (
	"log"
	"testing"
	"time"

	_ "github.com/fcg-xvii/go-tools/json/jsonmap"
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
				//go cl.Start()
			}
		}
	})
	go cl.Start()

	for {
		log.Println("Originate...")
		rOrig := &OriginateRequest{
			Channel: "SIP/user1/89774708408",
			Context: "admefine-bot",
			Exten:   "s",
		}
		originate, err := cl.Originate(rOrig)
		if err != nil {
			log.Println("Originate error", err)
			continue
		}
		log.Println("Call start")
		for event := range originate.Events() {
			log.Println(event)
		}
		log.Println("Call finished...")
		//time.Sleep(time.Second * 5)
		cl.Close()
	}
}
