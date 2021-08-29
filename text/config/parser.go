package config

import (
	"fmt"
	"io"
	"os"
)

type ParseMethod func(r io.Reader) (Config, error)

var (
	methods = make(map[string]ParseMethod)
)

func RegisterParseMethod(name string, method ParseMethod) {
	methods[name] = method
}

func FromReader(methodName string, r io.Reader) (res Config, err error) {
	method, check := methods[methodName]
	if !check {
		return nil, fmt.Errorf("UNEXPECTED METHOD [%v]", methodName)
	}
	return method(r)
}

func FromFile(methodName, filePath string) (res Config, err error) {
	method, check := methods[methodName]
	if !check {
		return nil, fmt.Errorf("UNEXPECTED METHOD [%v]", methodName)
	}
	var f *os.File
	if f, err = os.Open(filePath); err != nil {
		return
	}
	defer f.Close()
	return method(f)
}
