package ini

import (
	"fmt"
	"io"

	"github.com/fcg-xvii/go-tools/text/config"
	"github.com/fcg-xvii/go-tools/value"
)

func newConfig() *Config {
	mainSection := []config.Section{
		newSection(),
	}
	return &Config{
		sections: map[string][]config.Section{
			"main": mainSection,
		},
	}
}

type Config struct {
	sections map[string][]config.Section
}

func (s *Config) Sections(name string) ([]config.Section, bool) {
	sections, check := s.sections[name]
	return sections, check
}

func (s *Config) Section(name string) (config.Section, bool) {
	if sections, check := s.Sections(name); check {
		return sections[0], true
	}
	return nil, false
}

func (s *Config) AppendSection(name string) config.Section {
	section := newSection()
	if sections, check := s.Sections(name); check {
		sections = append(sections, section)
		s.sections[name] = sections
	} else {
		s.sections[name] = []config.Section{section}
	}
	return section
}

func (s *Config) ValueSetup(name string, ptr interface{}) bool {
	if main, check := s.Section("main"); check {
		return main.ValueSetup(name, ptr)
	}
	return false
}

func (s *Config) Value(name string) (value.Value, bool) {
	if main, check := s.Section("main"); check {
		return main.Value(name)
	}
	return value.Value{}, false
}

func (s *Config) ValueDefault(name string, defaultVal interface{}) interface{} {
	s.ValueSetup(name, &defaultVal)
	return defaultVal
}

func (s *Config) Save(w io.Writer) (err error) {
	for name, sections := range s.sections {
		for _, section := range sections {
			if _, err = w.Write([]byte(fmt.Sprintf("[%v]\n", name))); err != nil {
				return
			}
			if err = section.Save(w); err != nil {
				return
			}
		}
	}
	return
}
