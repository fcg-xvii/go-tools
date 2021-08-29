package config

import (
	"io"

	"github.com/fcg-xvii/go-tools/value"
)

type Section interface {
	ValueSetup(name string, ptr interface{}) bool
	Value(name string) (value.Value, bool)
	SetValue(name string, value interface{})
	Save(io.Writer) error
}

type Config interface {
	AppendSection(name string) (newSection Section)
	Section(name string) (Section, bool)
	Sections(name string) ([]Section, bool)
	ValueSetup(name string, ptr interface{}) bool
	Value(name string) (value.Value, bool)
	Save(io.Writer) error
}
