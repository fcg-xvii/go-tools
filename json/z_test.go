package json

import (
	"errors"
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

func (s *TimeInterval) JSONDecode(dec *JSONDecoder) (err error) {
	var sl []int
	if err = dec.Decode(&sl); err == nil {
		if len(sl) != 2 {
			return errors.New("TimeInterval: expected 2 int values")
		}
		s.Start, s.Finish = sl[0], sl[1]
	}
	return
}

type TimeIntervals []TimeInterval

func (s *TimeIntervals) JSONDecode(dec *JSONDecoder) (isNil bool, err error) {
	var sl [][]int
	if err = dec.Decode(&sl); err == nil && sl != nil {
		*s = make(TimeIntervals, 0, len(sl))
		for i, interval := range sl {
			if len(interval) != 2 {
				return true, fmt.Errorf("Index %v - expected 2 elements list", i)
			}
			*s = append(*s, TimeInterval{
				Start:  interval[0],
				Finish: interval[1],
			})
		}
	}
	isNil = sl == nil
	return
}

type NestedObject struct {
	ID        int
	Name      string
	Embedded  *TObject
	Intervals TimeIntervals
}

type TObject struct {
	ID        int
	Name      string
	Embedded  *NestedObject
	Intervals TimeIntervals
}

func (s *TObject) JSONField(fieldName string) (ptr interface{}, err error) {
	switch fieldName {
	case "id":
		ptr = &s.ID
	case "name":
		ptr = &s.Name
	case "embedded":
		ptr = &s.Embedded
	case "intervals":
		ptr = &s.Intervals
	}
	return
}

func TestDecoder(t *testing.T) {
	/*
		//src := []byte("[ [ 1, 2 ], [ 3, 4 ] ]")
		src := []byte("null")
		var sl *TimeIntervals
		err := DecodeBytes(src, &sl)
		log.Println(sl, err)
	*/
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
	log.Println("OBJ", obj, err, obj.Embedded, obj.Intervals)
	//log.Println("OBJ", obj, obj.embedded, obj.intervals)
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
	for _, v := range arr {
		log.Println(v)
	}
}
