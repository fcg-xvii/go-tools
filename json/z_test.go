package json

import (
	"log"
	"os"
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

type TObject struct {
	id       int
	name     string
	embedded *TObject
}

func (s *TObject) DecodeJSON(dec *JSONDecoder) error {
	return dec.DecodeObject(func(field string) (ptr interface{}, err error) {
		log.Println("<<<<<<<<<", field)
		switch field {
		case "id":
			ptr = &s.id
		case "name":
			ptr = &s.name
		case "embedded":
			ptr = &s.embedded
		}
		return
	})
}

func TestDecoder(t *testing.T) {
	// object
	fObj, err := os.Open("test_object.json")
	if err != nil {
		t.Error(err)
	}

	var obj TObject
	if err := Decode(fObj, &obj); err != nil {
		t.Error(err)
	}
	fObj.Close()
	log.Println("OBJ", obj, obj.embedded)
	// slice
	fObj, err = os.Open("test_array.json")
	if err != nil {
		t.Error(err)
	}

	var arr []*TObject
	if err := Decode(fObj, &arr); err != nil {
		t.Error(err)
	}
	fObj.Close()
	log.Println("ARR", arr)
}
