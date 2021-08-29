package ini

import (
	"bufio"
	"io"
	"strings"

	"github.com/fcg-xvii/go-tools/text/config"
)

func init() {
	config.RegisterParseMethod("ini", parser)
}

func parser(r io.Reader) (res config.Config, err error) {
	res = newConfig()
	buf := bufio.NewReader(r)
	var line []byte
	var section config.Section
	for {
		line, _, err = buf.ReadLine()
		if s := string(line); len(s) > 0 {
			s = strings.TrimSpace(s)
			switch s[0] {
			case '#':
				// comment
			case '[':
				// section
				if s[len(s)-1] == ']' {
					sectionName := s[1 : len(s)-1]
					if sectionName == "main" {
						section, _ = res.Section("main")
					} else {
						section = res.AppendSection(sectionName)
					}
				}
			default:
				// value, check comment
				if pos := strings.Index(s, "#"); pos > 0 {
					s = s[:pos]
				}
				// check plitter position
				if pos := strings.Index(s, "="); pos > 0 {
					key, val := strings.TrimSpace(s[:pos]), strings.TrimSpace(s[pos+1:])
					section.SetValue(key, val)
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
	}
}
