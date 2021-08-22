package ini

import (
	"os"
	"testing"

	"github.com/fcg-xvii/go-tools/text/config"
)

func TestINI(t *testing.T) {
	f, err := os.Open("test.ini")
	if err != nil {
		t.Fatal(err)
	}
	conf, err := config.FromReader("ini", f)
	t.Log(conf, err)
}
