package ini

import (
	"github.com/fcg-xvii/go-tools/value"
)

func newSection() *Section {
	return &Section{
		store: make(map[string]*value.Value),
	}
}

type Section struct {
	store map[string]*value.Value
}

func (s *Section) Value(name string, ptr interface{}) bool {
	if val, check := s.store[name]; check {
		return val.Setup(ptr)
	}
	return false
}
