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

	// read config
	conf, err := config.FromReader("ini", f)
	t.Log(conf, err)
	f.Close()

	// get value from config main section
	mainVal, check := conf.Value("one")
	t.Log(mainVal, check)
	// int
	var i int = 0
	t.Log(mainVal.Setup(&i), i)
	i = 0
	t.Log(conf.ValueSetup("one", &i), i)

	// get section
	main, check := conf.Section("main")
	t.Log(main, check)

	cool, check := conf.Section("cool")
	t.Log(cool, check)

	var str string = ""
	t.Log(cool.ValueSetup("key1", &str), str)

	cools, check := conf.Sections("cool")
	t.Log(cools, check)

	f, err = os.OpenFile("tmp.ini", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	conf.Save(f)
	f.Close()
}
