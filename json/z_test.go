package json

import (
	"log"
	"testing"
)

func TestJSON(t *testing.T) {
	m := Map{
		"jsrc": []int{1, 2, 3, 4},
		"kate": nil,
		"m1":   Map{"one": 1},
		"m2":   map[string]interface{}{"one": 2},
	}
	t.Log(m, string(m.JSON()))
	t.Log(string(m.ValueJSON("jsrc", []byte{})))
	t.Log(string(m.ValueJSON("jsrc1", []byte("{ }"))))
	t.Log(m.Map("m1", Map{}))
	t.Log(m.Map("m2", Map{}))
}

func TestInterface(t *testing.T) {
	m := MapFromInterface(Map{"one": 1})
	log.Println(m)
}
