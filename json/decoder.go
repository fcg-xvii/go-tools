package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type JSONObject interface {
	DecodeJSON(*JSONDecoder) error
}

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
	rVal := reflect.ValueOf(v)
	if rVal.Kind() == reflect.Ptr {
		rVal = rVal.Elem()
		if rVal.Kind() == reflect.Ptr && rVal.IsNil() {
			rVal.Set(reflect.New(rVal.Type().Elem()))
			v = rVal.Elem().Addr().Interface()
		}
		s.parentObj = rVal
	}
	if jsonObj, check := v.(JSONObject); check {
		return jsonObj.DecodeJSON(s)
	} else {
		if rVal.Kind() == reflect.Slice {
			err = s.decodeSlice(&rVal)
		} else {
			err = s.Decoder.Decode(v)
		}
	}
	if err != nil {
		s.err = err
	}
	return
}

func (s *JSONDecoder) decodeSlice(sl *reflect.Value) (err error) {
	if _, err = s.Token(); err != nil {
		return
	}
	if s.current != JSON_ARRAY {
		if s.token == nil {
			if s.parentObj.IsValid() {
				s.parentObj.Set(reflect.Zero(s.parentObj.Type()))
				s.parentObj = reflect.Value{}
			}
			return nil
		} else {
			return fmt.Errorf("EXPECTED ARRAY, NOT %T", s.current)
		}
	}
	elemType := reflect.TypeOf(sl.Interface()).Elem()
	for s.More() {
		if !s.More() {
			return nil
		}
		var eType reflect.Type
		if elemType.Kind() == reflect.Slice {
			eType = elemType
		} else if elemType.Kind() == reflect.Ptr {
			eType = reflect.PtrTo(elemType.Elem())
		} else {
			eType = reflect.PtrTo(elemType)
		}
		rElem := reflect.New(eType)
		val := rElem.Interface()
		if err = s.Decode(val); err != nil {
			return
		}
		if !s.parentObj.IsValid() {
			sl.Set(reflect.Append(*sl, reflect.Zero(elemType)))
		} else {
			if elemType.Kind() == reflect.Ptr || elemType.Kind() == reflect.Slice {
				sl.Set(reflect.Append(*sl, rElem.Elem()))
			} else {
				sl.Set(reflect.Append(*sl, rElem.Elem().Elem()))
			}
		}
	}
	if _, err = s.Token(); err != nil {
		return
	}
	if d, check := s.token.(json.Delim); !check || d != ']' {
		return fmt.Errorf("JSON PARSE ERROR :: EXPECTED ']', NOT %v", d)
	}
	return nil
}

func (s *JSONDecoder) DecodeObject(fieldRequest func(string) (interface{}, error)) (err error) {
	if _, err = s.Token(); err != nil {
		return
	}
	if s.current != JSON_OBJECT {
		if s.token == nil {
			if s.parentObj.IsValid() {
				s.parentObj.Set(reflect.Zero(s.parentObj.Type()))
				s.parentObj = reflect.Value{}
			}
			return
		} else {
			return fmt.Errorf("EXPCTED OBJECT, NOT %T", s.current)
		}
	}
	el := s.EmbeddedLevel()
	var t json.Token
	var ptr interface{}
	for el <= s.EmbeddedLevel() {
		if t, err = s.Token(); err != nil {
			return
		}
		if s.Current() == JSON_VALUE && s.IsObjectKey() {
			if ptr, err = fieldRequest(t.(string)); err != nil {
				return
			}
			if ptr != nil {
				err = s.Decode(ptr) // err = nil (WTF???)
				if s.err != nil {
					err = s.err
					return
				}
			} else {
				if err = s.Next(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
