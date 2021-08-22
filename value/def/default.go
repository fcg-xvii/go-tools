package def

import (
	"strconv"
	"strings"
)

func String(val, defaultVal string) (res string) {
	if len(val) == 0 {
		res = defaultVal
	}
	return
}

func StringTrim(val, defaultVal string) string {
	return String(strings.TrimSpace(val), defaultVal)
}

func Int(val string, defaultVal int) (res int) {
	var err error
	if res, err = strconv.Atoi(val); err != nil {
		res = defaultVal
	}
	return
}
