package jsonmap

import (
	"testing"
)

func TestJSON(t *testing.T) {
	m := JSONMap{
		"jsrc": []int{1, 2, 3, 4},
	}

	t.Log(string(m.ValueJSON("jsrc", []byte{})))
	t.Log(string(m.ValueJSON("jsrc1", []byte("{ }"))))
}
