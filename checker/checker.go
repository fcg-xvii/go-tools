package checker

import (
	"regexp"
	"strings"
)

type Type byte

const (
	Undefined Type = iota
	Email
	Phone
)

func (s Type) String() string {
	switch s {
	case Email:
		return "email"
	case Phone:
		return "phone"
	default:
		return "undefined"
	}
}

var (
	checkers = map[Type]*regexp.Regexp{
		Email: regexp.MustCompile("^.*?@.*?\\.[\\w]+$"),
		Phone: regexp.MustCompile("^\\+\\d{11}$"),
	}
)

func Check(source string) Type {
	source = strings.TrimSpace(source)
	for i, v := range checkers {
		if v.MatchString(source) {
			return i
		}
	}
	return Undefined
}

func CheckEmail(email string) bool {
	return checkers[Email].MatchString(email)
}

func CheckPhone(phone string) bool {
	return checkers[Phone].MatchString(phone)
}
