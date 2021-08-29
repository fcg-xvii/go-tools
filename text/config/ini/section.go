package ini

import (
	"fmt"
	"io"

	"github.com/fcg-xvii/go-tools/text/config"
	"github.com/fcg-xvii/go-tools/value"
)

func newSection() config.Section {
	return &Section{
		store: make(map[string]value.Value),
	}
}

type Section struct {
	store map[string]value.Value
}

func (s *Section) ValueSetup(name string, ptr interface{}) bool {
	if val, check := s.store[name]; check {
		return val.Setup(ptr)
	}
	return false
}

func (s *Section) Value(name string) (value.Value, bool) {
	val, check := s.store[name]
	return val, check
}

func (s *Section) SetValue(name string, val interface{}) {
	s.store[name] = value.ValueOf(val)
}

func (s *Section) Save(w io.Writer) (err error) {
	for name, val := range s.store {
		if _, err = w.Write([]byte(fmt.Sprintf("%v = %v\n", name, val.String()))); err != nil {
			return err
		}
	}
	_, err = w.Write([]byte("\n"))
	return
}
