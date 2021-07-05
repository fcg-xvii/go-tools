package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"

	"github.com/fcg-xvii/go-tools/containers"
)

type JSONTokenType byte

const (
	JSON_INVALID JSONTokenType = iota
	JSON_ARRAY
	JSON_OBJECT
	JSON_VALUE
)

func (s JSONTokenType) String() string {
	switch s {
	case JSON_INVALID:
		return "JSON_INVALID"
	case JSON_ARRAY:
		return "JSON_ARRAY"
	case JSON_OBJECT:
		return "JSON_OBJECT"
	case JSON_VALUE:
		return "JSON_VALUE"
	default:
		return "JSON_UNDEFINED"
	}
}

func Decode(r io.Reader, obj interface{}) error {
	dec := InitJSONDecoder(r)
	return dec.Decode(obj)
}

func DecodeBytes(src []byte, obj interface{}) error {
	buf := bytes.NewBuffer(src)
	dec := InitJSONDecoder(buf)
	return dec.Decode(obj)
}

// JSON decode interfaces

type Type byte

const (
	TypeUndefined Type = iota
	TypeInterface
	TypeObject
	TypeSlice
)

type JSONInterface interface {
	JSONDecode(*JSONDecoder) (isNil bool, err error)
}

type JSONObject interface {
	JSONField(fieldName string) (fieldPtr interface{}, err error)
}

func (s *JSONDecoder) JSONTypeCheck(rv *reflect.Value) (t Type) {
	// check slice
	if rv.Kind() == reflect.Slice {
		return TypeSlice
	}
	iface := rv.Interface()
	if _, check := iface.(JSONInterface); check {
		// check JSONInterface
		t = TypeInterface
		return
	} else if _, check = iface.(JSONObject); check {
		// check JSONObject
		t = TypeObject
		return
	}
	return
}

///////////////////////////////////////////////

func InitJSONDecoderFromSource(src []byte) *JSONDecoder {
	r := bytes.NewReader(src)
	return InitJSONDecoder(r)
}

func InitJSONDecoder(r io.Reader) *JSONDecoder {
	return &JSONDecoder{
		Decoder:  json.NewDecoder(r),
		embedded: containers.NewStack(0),
	}
}

type JSONDecoder struct {
	*json.Decoder
	token        json.Token
	embedded     *containers.Stack
	current      JSONTokenType
	objectkey    bool
	objectClosed bool
	parentObj    reflect.Value
	err          error
	counter      int
}

func (s *JSONDecoder) IsObjectKey() bool      { return s.objectkey }
func (s *JSONDecoder) IsObjectClosed() bool   { return s.objectClosed }
func (s *JSONDecoder) Current() JSONTokenType { return s.current }
func (s *JSONDecoder) EmbeddedLevel() int     { return s.embedded.Len() }

func (s *JSONDecoder) Token() (t json.Token, err error) {
	s.objectClosed = false
	if t, err = s.Decoder.Token(); err == nil {
		if delim, check := t.(json.Delim); check {
			s.objectkey = false
			switch delim {
			case '{':
				s.embedded.Push(JSON_OBJECT)
				s.current = JSON_OBJECT
			case '[':
				s.embedded.Push(JSON_ARRAY)
				s.current = JSON_ARRAY
			case '}', ']':
				s.embedded.Pop()
				s.objectClosed, s.current = true, JSON_INVALID
				if s.embedded.Len() > 0 {
					s.current = s.embedded.Peek().(JSONTokenType)
				}
			}
		} else {
			if s.current == JSON_OBJECT {
				s.objectkey = !s.objectkey
			}
			s.current = JSON_VALUE
		}
	}
	s.token = t
	return
}

func (s *JSONDecoder) Next() error {
	if _, err := s.Token(); err != nil {
		return err
	}
	switch s.current {
	case JSON_ARRAY, JSON_OBJECT:
		{
			stackLen := s.embedded.Len()
			for s.embedded.Len() >= stackLen {
				if _, err := s.Token(); err != nil {
					return err
				}
			}
			return nil
		}
	default:
		return nil
	}
}

func (s *JSONDecoder) DecodeRaw(v interface{}) error {
	return s.Decoder.Decode(v)
}

func (s *JSONDecoder) Decode(v interface{}) (err error) {
	rv := reflect.ValueOf(v)
	err = s.decodeReflect(&rv)
	return
}

func (s *JSONDecoder) decodeReflect(rv *reflect.Value) (err error) {
	switch s.JSONTypeCheck(rv) {
	case TypeInterface:
		{
			if rv.Kind() == reflect.Ptr && rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			var isNil bool
			if isNil, err = rv.Interface().(JSONInterface).JSONDecode(s); err == nil {
				if isNil && rv.CanAddr() {
					rv.Set(reflect.Zero(rv.Type()))
				}
			}
			return
		}
	case TypeObject:
		return s.decodeJSONObject(rv)
	case TypeSlice:
		return s.decodeSlice(rv)
	default:
		{
			if rv.Kind() == reflect.Ptr {
				ev := rv.Elem()
				if !ev.IsValid() {
					ev = reflect.New(rv.Type().Elem())
					//log.Println(ev.Elem().Kind(), ev.Elem().IsNil())
					if ev.Elem().Kind() == reflect.Ptr {
						if err = s.decodeReflect(&ev); err == nil {
							rv.Set(ev)
						}
						return
					} else {
						if err = s.Decoder.Decode(rv.Interface()); err == nil {
							rv.Set(ev)
							log.Println(">>>", ev.Elem(), ev.Interface() == nil)
						}
						return
					}
				} else {
					return s.decodeReflect(&ev)
				}
			} else {
				pv := reflect.New(rv.Type())
				if err = s.Decoder.Decode(pv.Interface()); err == nil {
					rv.Set(pv.Elem())
				}
			}
			return
		}
	}
}

func (s *JSONDecoder) decodeRawObject(rv *reflect.Value) (err error) {
	var t json.Token
	// chek first token is object
	if t, err = s.Token(); err != nil {
		return
	}
	if s.current != JSON_OBJECT {
		// check null token
		if t == nil {
			if rv.CanAddr() && !rv.IsNil() {
				rv.Set(reflect.Zero(rv.Type()))
			}
			return
		}
		// token is not object
		return fmt.Errorf("EXPCTED OBJECT, NOT %T", t)
	}
	// check null pounter in source object
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		// create new object
		rv.Set(reflect.New(rv.Type().Elem()))
	}
	el := s.EmbeddedLevel()
	for el <= s.EmbeddedLevel() {
		if t, err = s.Token(); err != nil {
			return
		}
		if s.Current() == JSON_VALUE && s.IsObjectKey() {
			log.Println(t.(string))

			/*
				if fieldPtr, err = obj.JSONField(t.(string)); err != nil {
					return
				}
				if fieldPtr != nil {
					rv := reflect.ValueOf(fieldPtr)
					if err = s.decodeReflect(&rv); err != nil {
						return
					}
				} else {
					if err = s.Next(); err != nil {
						return
					}
				}
			*/
		}
	}
	return
}

// slice
func (s *JSONDecoder) decodeSlice(rv *reflect.Value) (err error) {
	var t json.Token
	if t, err = s.Token(); err != nil {
		return
	}
	if s.current != JSON_ARRAY {
		if t == nil {
			if !rv.IsNil() {
				rv.Set(reflect.Zero(rv.Type()))
			}
			return
		}
		return fmt.Errorf("EXPCTED SLICE, NOT %T", t)
	}
	// check slice is nil
	if rv.IsNil() {
		rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
	}
	elemType := reflect.TypeOf(rv.Interface()).Elem()
	for s.More() {
		em := reflect.New(elemType)
		if err = s.decodeReflect(&em); err != nil {
			return
		}
		rv.Set(reflect.Append(*rv, em.Elem()))
	}
	if _, err = s.Token(); err != nil {
		return
	}
	if d, check := s.token.(json.Delim); !check || d != ']' {
		return fmt.Errorf("JSON PARSE ERROR :: EXPECTED ']', NOT %v", d)
	}
	return
}

////////////////////////////////////////

func (s *JSONDecoder) decodeJSONObject(rv *reflect.Value) (err error) {
	var t json.Token
	// chek first token is object
	if t, err = s.Token(); err != nil {
		return
	}
	if s.current != JSON_OBJECT {
		// check null token
		if t == nil {
			if rv.CanAddr() && !rv.IsNil() {
				rv.Set(reflect.Zero(rv.Type()))
			}
			return
		}
		// token is not object
		return fmt.Errorf("EXPCTED OBJECT, NOT %T", t)
	}
	// check null pounter in source object
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		// create new object
		rv.Set(reflect.New(rv.Type().Elem()))
	}
	obj := rv.Interface().(JSONObject)
	el := s.EmbeddedLevel()
	var fieldPtr interface{}
	for el <= s.EmbeddedLevel() {
		if t, err = s.Token(); err != nil {
			return
		}
		if s.Current() == JSON_VALUE && s.IsObjectKey() {
			if fieldPtr, err = obj.JSONField(t.(string)); err != nil {
				return
			}
			if fieldPtr != nil {
				rv := reflect.ValueOf(fieldPtr)
				if err = s.decodeReflect(&rv); err != nil {
					return
				}
			} else {
				if err = s.Next(); err != nil {
					return
				}
			}
		}
	}
	return
}
