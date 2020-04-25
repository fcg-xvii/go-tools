package ami

import (
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	cl := Open("127.0.0.1:5538", "admin", "eternalsun")

	cl.Request(Action{
		"Action": "Originate",
	})
	time.Sleep(time.Second * 30)
	cl.Close()
}
