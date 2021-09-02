package value

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ValueOf(val interface{}) Value {
	return Value{
		val: val,
	}
}

type Value struct {
	val interface{}
}

func (s *Value) IsValid() bool {
	return s.val != nil
}

func (s *Value) Setup(val interface{}) (res bool) {
	rr := reflect.ValueOf(val)
	if rr.Kind() != reflect.Ptr {
		panic(fmt.Errorf("Expected ptr, given %v", rr.Type()))
	}
	rKind, rType := rr.Elem().Kind(), rr.Elem().Type()
	if rKind == reflect.Interface {
		rKind, rType = rr.Elem().Elem().Kind(), rr.Elem().Elem().Type()
	}
	if rKind == reflect.String {
		rls := strings.TrimSpace(fmt.Sprint(s.val))
		if len(rls) > 0 {
			rr.Elem().Set(reflect.ValueOf(rls))
			res = true
		}
	} else {
		rl := reflect.ValueOf(s.val)
		if rl.Kind() == reflect.String {
			switch rKind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				{
					if tmp, err := strconv.ParseInt(rl.String(), 10, 64); err == nil {
						rr.Elem().Set(reflect.ValueOf(tmp).Convert(rType))
						res = true
					}
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				{
					if tmp, err := strconv.ParseUint(rl.String(), 10, 64); err == nil {
						rr.Elem().Set(reflect.ValueOf(tmp).Convert(rType))
						res = true
					}
				}
			case reflect.Float32, reflect.Float64:
				{
					if tmp, err := strconv.ParseFloat(rl.String(), 64); err == nil {
						rr.Elem().Set(reflect.ValueOf(tmp).Convert(rType))
						res = true
					}
				}
			default:
				// json
				i := reflect.New(rr.Elem().Type()).Interface()
				if err := json.Unmarshal([]byte(rl.String()), i); err == nil {
					rr.Elem().Set(reflect.ValueOf(i).Elem())
				}
			}
		} else {
			var rVal reflect.Value
			defer func() {
				if r := recover(); r == nil {
					rr.Elem().Set(rVal)
					res = true
				}
			}()
			rVal = rl.Convert(rType)
		}
	}
	return
}

func (s *Value) String() string {
	return fmt.Sprint(s.val)
}

func (s *Value) Int() int {
	var i int
	s.Setup(&i)
	return i
}

func (s *Value) Int8() int8 {
	var i int8
	s.Setup(&i)
	return i
}

func (s *Value) Int16() int16 {
	var i int16
	s.Setup(&i)
	return i
}

func (s *Value) Int32() int32 {
	return int32(s.Int())
}

func (s *Value) Float32() float32 {
	var i float32
	s.Setup(&i)
	return i
}

func (s *Value) Float64() float64 {
	var i float64
	s.Setup(&i)
	return i
}
