package value

import (
	"log"
	"testing"
)

func TestValue(t *testing.T) {
	val := ValueOf(`{ "one": "111" }`)
	s := map[string]string{}
	t.Log(val.Setup(&s), s)

	val = ValueOf(100.55)
	log.Println(val.Int())
}
