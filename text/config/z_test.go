package config

import (
	"testing"
)

var (
	splitFile = "z_split.config"
)

func TestSplitFile(t *testing.T) {
	if s, err := SplitFile(splitFile, "::"); err != nil {
		t.Error(err)
	} else {
		t.Log(s)
	}
}

func TestSplitFileMap(t *testing.T) {
	if m, err := SplitFileToMap(splitFile, "::", "host", "login", "password", "undefined"); err != nil {
		t.Error(err)
	} else {
		t.Log(m)
	}
}

func TestSplitFileVals(t *testing.T) {
	var host, login, password string
	if err := SplitFileToVals(splitFile, "::", &host, &login, &password); err != nil {
		t.Error(err)
	}
	t.Log(host, login, password)
}
