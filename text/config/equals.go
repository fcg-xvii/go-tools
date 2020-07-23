package config

// Map with Get method
type Map map[string]string

// Get returns value by key or defaultVal if key is not exists
func (s Map) Get(key, defaultVal string) (res string) {
	var check bool
	if res, check = s[key]; !check {
		res = defaultVal
	}
	return
}
