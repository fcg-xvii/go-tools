package json

import (
	"fmt"
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

type TimeInterval struct {
	Start  int
	Finish int
}

type TimeIntervals []TimeInterval

func (s *TimeIntervals) JSONDecode(dec *JSONDecoder) error {
	var sl [][]int
	err := dec.Decode(&sl)
	if err == nil {
		log.Println("SSSSSS", s)
		*s = make(TimeIntervals, 0, len(sl))
		for i, interval := range sl {
			if len(interval) != 2 {
				return fmt.Errorf("Index %v - expected 2 elements list", i)
			}
			*s = append(*s, TimeInterval{
				Start:  interval[0],
				Finish: interval[1],
			})
		}
	}
	return err
}

type TObject struct {
	id        int
	name      string
	embedded  *TObject
	intervals *TimeIntervals
}

func (s *TObject) JSONField(fieldName string) (ptr interface{}, err error) {
	switch fieldName {
	case "id":
		ptr = &s.id
	case "name":
		ptr = &s.name
	case "embedded":
		ptr = &s.embedded
	case "intervals":
		log.Println("DDD", s.intervals)
		ptr = &s.intervals
	}
	return
	/*
		return dec.DecodeObject(func(field string) (ptr interface{}, err error) {
			log.Println("<<<<<<<<<", field)
			switch field {
			case "id":
				ptr = &s.id
			case "name":
				ptr = &s.name
			case "embedded":
				ptr = &s.embedded
			case "intervals":
				log.Println("DDD", s.intervals)
				ptr = &s.intervals
			}
			return
		})
	*/
}

/*func (s *TObject) EmbeddedString() string {
	if s.embedded == nil {
		return "nil"
	} else {
		return s.embedded.String()
	}
}

func (s *TObject) String() string {
	m := map[string]interface{}{
		"id":        s.id,
		"name":      s.name,
		"embedded":  s.EmbeddedString(),
		"intervals": s.intervals,
	}
	str, _ := json.MarshalIndent(m, "", "\t")
	return string(str)
}*/

func TestDecoder(t *testing.T) {
	/*
		src := []byte("[ [ 1, 2 ], [ 3, 4 ] ]")
		var sl [][]int
		err := DecodeBytes(src, &sl)
		log.Println(sl, err)
	*/
	// object
	fObj, err := os.Open("test_object.json")
	if err != nil {
		t.Error(err)
	}

	var obj *TObject = new(TObject)
	if err := Decode(fObj, &obj); err != nil {
		t.Error(err)
	}
	fObj.Close()
	log.Println("OBJ", obj, err)
	//log.Println("OBJ", obj, obj.embedded, obj.intervals)
	// slice
	/*
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
	*/
}
