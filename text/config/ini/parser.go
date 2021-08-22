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
	buf := bufio.NewReader(r)
	var line []byte
	var section *Section
	for {
		line, _, err = buf.ReadLine()
		if s := string(line); len(s) > 0 {
			switch s[0] {
			case '#':
				// comment
			case '[':
				// section
				s = strings.TrimSpace(s)
				if s[len(s)-1:] != ']'
					break
			default:
				// value
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
