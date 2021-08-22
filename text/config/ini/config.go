package ini

func newConfig() *Config {
	return &Config{
		sections: map[string][]*Section{
			"main": []*Sections{
				newSection(),
			},
		}
	}
}

type Config struct {
	sections map[string][]*Section
}

func (s *Config) Sections(name string) ([]*Section, bool) {
	sections, check := s.sections[name]
	return sections, check
}

func (s *Config) Section(name string) (*Section, bool) {
	if sections, check := s.Sections(name); check {
		return sections[0], true
	}
	return nil, false
}
