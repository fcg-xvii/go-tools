package config

type Section interface {
	Value(name string, ptr interface{}) bool
}

type Config interface {
	Section(name string) Section
	Value(name string, ptr interface{}) bool
}
