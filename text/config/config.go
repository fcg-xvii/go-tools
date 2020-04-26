package config

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// SplitFile read data from file and returns a slice of lines from a read stream and split by splitter argument
func SplitFile(fileName string, splitter string) (res []string, err error) {
	var f *os.File
	if f, err = os.Open(fileName); err != nil {
		return
	}
	res, err = Split(f, splitter)
	f.Close()
	return
}

// Split returns a slice of lines from a read stream and split by splitter argument
func Split(r io.Reader, splitter string) (res []string, err error) {
	var source []byte
	if source, err = ioutil.ReadAll(r); err != nil {
		return
	}
	res = strings.Split(string(source), splitter)
	return
}

func partsToMap(parts []string, keys []string) map[string]string {
	res := make(map[string]string)
	for i, key := range keys {
		if i < len(parts) {
			res[key] = parts[i]
		} else {
			res[key] = ""
		}
	}
	return res
}

func partsToVals(parts []string, ptrs []*string) {
	for i, ptr := range ptrs {
		if i < len(parts) {
			*ptr = parts[i]
		} else {
			*ptr = ""
		}
	}
}

//
func SplitToMap(r io.Reader, splitter string, keys ...string) (res map[string]string, err error) {
	var parts []string
	if parts, err = Split(r, splitter); err != nil {
		return
	}
	res = partsToMap(parts, keys)
	return
}

//
func SplitFileToMap(fileName string, splitter string, keys ...string) (res map[string]string, err error) {
	var parts []string
	if parts, err = SplitFile(fileName, splitter); err != nil {
		return
	}
	res = partsToMap(parts, keys)
	return
}

func SplitToVals(r io.Reader, splitter string, ptrs ...*string) (err error) {
	var parts []string
	if parts, err = Split(r, splitter); err != nil {
		return
	}
	partsToVals(parts, ptrs)
	return
}

func SplitFileToVals(fileName string, splitter string, ptrs ...*string) (err error) {
	var parts []string
	if parts, err = SplitFile(fileName, splitter); err != nil {
		return
	}
	partsToVals(parts, ptrs)
	return
}
