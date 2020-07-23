package config

import (
	"os"
	"testing"
)

var (
	splitFile = "z_split.config"
)

func checkConfigExists() bool {
	_, err := os.Stat(splitFile)
	return err == nil
}

func TestSplitFile(t *testing.T) {
	if !checkConfigExists() {
		t.Log("Expected config file ", splitFile)
		return
	}
	if s, err := SplitFile(splitFile, "::"); err != nil {
		t.Error(err)
	} else {
		t.Log(s)
	}
}

func TestSplitFileMap(t *testing.T) {
	if !checkConfigExists() {
		t.Log("Expected config file ", splitFile)
		return
	}
	if m, err := SplitFileToMap(splitFile, "::", "host", "login", "password", "undefined"); err != nil {
		t.Error(err)
	} else {
		t.Log(m)
	}
}

func TestSplitFileVals(t *testing.T) {
	if !checkConfigExists() {
		t.Log("Expected config file ", splitFile)
		return
	}
	var host, login, password string
	if err := SplitFileToVals(splitFile, "::", &host, &login, &password); err != nil {
		t.Error(err)
	}
	t.Log(host, login, password)
}

func TestEqual(t *testing.T) {
	m, err := SplitFileToEquals("test_content/equal_conf.cfg")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(m)
}
