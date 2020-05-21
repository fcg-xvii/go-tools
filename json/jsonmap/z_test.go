package jsonmap

import (
	"testing"
)

func TestJSON(t *testing.T) {
	m := JSONMap{
		"jsrc": []int{1, 2, 3, 4},
		"kate": nil,
	}
	t.Log(m, string(m.JSON()))
	t.Log(string(m.ValueJSON("jsrc", []byte{})))
	t.Log(string(m.ValueJSON("jsrc1", []byte("{ }"))))
}
