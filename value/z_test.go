package value

import "testing"

func TestValue(t *testing.T) {
	val := ValueOf(`{ "one": "111" }`)
	s := map[string]string{}
	t.Log(val.Setup(&s), s)
}
