package jsonmap

import (
	"testing"
)

func TestJSON(t *testing.T) {
	m := JSONMap{
		"jsrc": []int{1, 2, 3, 4},
		"kate": nil,
		"m1":   JSONMap{"one": 1},
		"m2":   map[string]interface{}{"one": 2},
	}
	t.Log(m, string(m.JSON()))
	t.Log(string(m.ValueJSON("jsrc", []byte{})))
	t.Log(string(m.ValueJSON("jsrc1", []byte("{ }"))))
	t.Log(m.JSONMap("m1", JSONMap{}))
	t.Log(m.JSONMap("m2", JSONMap{}))
}
