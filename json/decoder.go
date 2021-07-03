package json

import (
	"bytes"
	"encoding/json"
	"errors"
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
	JSONDecode(*JSONDecoder) error
}

type JSONObject interface {
	JSONField(fieldName string) (fieldPtr interface{}, err error)
}

func (s *JSONDecoder) JSONTypeCheck(rv *reflect.Value) (t Type) {
	// check slice
	if rv.Kind() == reflect.Slice {
		return TypeSlice
	}
	// check to get interface
	if !rv.CanInterface() {
		return
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
	log.Println("DECODE...")
	rv := reflect.ValueOf(v)
	err = s.decodeReflect(&rv)
	log.Println("UUU", err)
	return
}

func (s *JSONDecoder) decodeReflect(rv *reflect.Value) (err error) {
	log.Println("DECODE_REFLECT", rv, rv.IsNil())
	switch s.JSONTypeCheck(rv) {
	case TypeInterface:
		{
			if rv.Kind() == reflect.Ptr && rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			log.Println("RR", rv, rv.IsNil())
			return rv.Interface().(JSONInterface).JSONDecode(s)
		}
	case TypeObject:
		return s.decodeJSONObject(rv)
	default:
		{
			if rv.Kind() == reflect.Ptr {
				ev := rv.Elem()
				if ev.Kind() == reflect.Ptr {
					if err = s.decodeReflect(&ev); err != nil {
						return
					}
					if rv.IsNil() {
						log.Println("NILLLLLL")
						return errors.New("NIL")
					}
					return
				}
			}
			if err = s.Decoder.Decode(rv.Interface()); err == nil {
				log.Println("IS_NIL...", rv.IsNil())
			}
			return
		}
	}
}

// slice

func (s *JSONDecoder) decodeSlice(rv *reflect.Value) (err error) {
	log.Println("decode slice")
	var t json.Token
	if t, err = s.Token(); err != nil {
		return
	}
	if s.current != JSON_ARRAY {
		log.Println("NOT ARRAY")
		if t == nil {
			if !rv.IsNil() {
				rv.Set(reflect.Zero(rv.Type()))
			}
			return
		}
		return fmt.Errorf("EXPCTED SLICE, NOT %T", t)
	}
	// check slice is nil
	log.Println(rv, rv.Type())
	if rv.IsNil() {
		//rv.Set(reflect.New(rv.Type()))
		rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
	}
	elemType := reflect.TypeOf(rv.Interface()).Elem()
	log.Println("ELEM_TYPE", elemType)
	for s.More() {
		em := reflect.New(elemType)
		if err = s.decodeReflect(&em); err != nil {
			return
		}
		//log.Println(em.Elem().Interface())
		rv.Set(reflect.Append(*rv, em.Elem()))
	}
	log.Println(rv.Interface())
	return
}

////////////////////////////////////////

func (s *JSONDecoder) decodeJSONObject(rv *reflect.Value) (err error) {
	log.Println("DECODE OBJ")
	var t json.Token
	// chek first token is object
	if t, err = s.Token(); err != nil {
		return
	}
	if s.current != JSON_OBJECT {
		log.Println("CURRENT", t == nil)
		// check null token
		if t == nil {
			log.Println(rv, rv.CanAddr(), rv.IsNil())
			if rv.CanAddr() && !rv.IsNil() {
				rv.Set(reflect.Zero(rv.Type()))
			}
			return
		}
		// token is not object
		return fmt.Errorf("EXPCTED OBJECT, NOT %T", t)
	}
	// check null pounter in source object
	if rv.IsNil() {
		// create new object
		rv.Set(reflect.New(rv.Type()))
		log.Println("new object created")
	}
	obj := rv.Interface().(JSONObject)
	el := s.EmbeddedLevel()
	var fieldPtr interface{}
	for el <= s.EmbeddedLevel() {
		if t, err = s.Token(); err != nil {
			return
		}
		if s.Current() == JSON_VALUE && s.IsObjectKey() {
			log.Println("FIELD_NAME", t.(string))
			if fieldPtr, err = obj.JSONField(t.(string)); err != nil {
				return
			}
			if fieldPtr != nil {
				rv := reflect.ValueOf(fieldPtr)
				if err = s.decodeReflect(&rv); err != nil {
					log.Println("ERRRRRRR", err)
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

/*
func (s *JSONDecoder) decodeSlice(sl *reflect.Value) (err error) {
	log.Println("DECODE_SLICE")
	if _, err = s.Token(); err != nil {
		return
	}
	if s.current != JSON_ARRAY {
		if s.token == nil {
			log.Println("AAAAAAAAAAAAAAAAAAAAAAAA")
			sl.Set(reflect.Zero(sl.Type()))
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
	log.Println("ELEM_TYPE")
	for s.More() {
		log.Println("!!!!!")
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
		log.Println("VALLLLL", val)
		if !s.parentObj.IsValid() {
			sl.Set(reflect.Append(*sl, reflect.Zero(elemType)))
		} else {
			if elemType.Kind() == reflect.Ptr || elemType.Kind() == reflect.Slice {
				log.Println(rElem.Elem().IsNil())
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

func (s *JSONDecoder) decodeObject(fieldRequest func(string) (interface{}, error)) (err error) {
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
*/
