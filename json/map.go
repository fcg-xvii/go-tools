package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Map type
type Map map[string]interface{}

// New init Map object
func NewMap() Map {
	return make(Map)
}

// FromInterface convert map[string]interface{} or Map interface to Map
func MapFromInterface(iface interface{}) (res Map) {
	switch val := iface.(type) {
	case map[string]interface{}:
		res = FromMap(val)
	case Map:
		res = val
	default:
		res = NewMap()
	}
	return
}

// FromMap convert map[string]interface{} to Map object
func FromMap(m map[string]interface{}) Map {
	return Map(m)
}

// Bool returns bool value by key
func (s Map) Bool(key string, defaultVal bool) bool {
	if res, check := s[key].(bool); check {
		return res
	}
	return defaultVal
}

// KeyExists check value exists by key
func (s Map) KeyExists(key string) bool {
	_, check := s[key]
	return check
}

// KeysExists check values exists by keys list.
// Returns the first blank key found.
// If all keys are defined, an empty string will be returned
func (s Map) KeysExists(keys []string) string {
	for _, key := range keys {
		if _, check := s[key]; !check {
			return key
		}
	}
	return ""
}

func val(l, r interface{}) (res reflect.Value) {
	lVal, rVal := reflect.ValueOf(l), reflect.ValueOf(r)
	if lVal.Kind() == reflect.Ptr && rVal.Kind() != reflect.Ptr {
		return val(lVal.Elem().Interface(), r)
	}
	defer func() {
		if r := recover(); r != nil {
			res = rVal
		}
	}()
	res = lVal.Convert(rVal.Type())
	return
}

// Int returns int64 value by key.
// If key isn't defined will be returned defaultVal arg value
func (s Map) Int(key string, defaultVal int64) int64 {
	if iface, check := s[key]; check {
		return val(iface, defaultVal).Int()
	}
	return defaultVal
}

func (s Map) Int32(key string, defaultVal int) int {
	if iface, check := s[key]; check {
		return val(iface, defaultVal).Interface().(int)
	}
	return defaultVal
}

// Value returns interface object with attempt to convert to defaultVal type.
// If key isn't defined will be returned defaultVal arg value
func (s Map) Value(key string, defaultVal interface{}) interface{} {

	// check value exists by key
	if iface, check := s[key]; check {

		// check defaultVal is valid
		dVal := reflect.ValueOf(defaultVal)
		if !dVal.IsValid() {
			// invalid, return arrived interface
			return iface
		}
		// defaultVal is valid, attempt to convert found value to defaultVal type
		lVal := reflect.ValueOf(iface)

		switch {
		case !lVal.IsValid():
			return defaultVal // invalid found value, return defaultVal
		case lVal.Kind() == dVal.Kind():
			return iface // types of found value and defaultVal is match. return found value
		case lVal.Type().ConvertibleTo(dVal.Type()):
			return lVal.Convert(dVal.Type()).Interface() // found value type can be converted to defaultVal type. return converted found value
		default:
			{
				// found value type can't be converted to defaultVal type. If found value is string, attempt to convert to number value
				if lVal.Kind() == reflect.String {
					switch dVal.Kind() {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						val, err := strconv.Atoi(lVal.String())
						if err != nil {
							return defaultVal
						}
						lVal = reflect.ValueOf(val)
						return lVal.Convert(dVal.Type())
					case reflect.Float32, reflect.Float64:
						val, err := strconv.ParseFloat(lVal.String(), 64)
						if err != nil {
							return defaultVal
						}
						lVal = reflect.ValueOf(val)
						return lVal.Convert(dVal.Type())
					}
				}
			}
		}
	}
	return defaultVal
}

// ValueJSON returns json source object by key
// If key isn't defined will be returned defaultVal arg value
func (s Map) ValueJSON(key string, defaultVal []byte) (res []byte) {
	res = defaultVal
	if s.KeyExists(key) {
		if r, err := json.Marshal(s[key]); err == nil {
			res = r
		}
	}
	return
}

// String returns string value by key
// If key isn't defined will be returned defaultVal arg value
func (s Map) String(key, defaultVal string) string {
	if iface, check := s[key]; check {
		return val(iface, defaultVal).String()
		//return fmt.Sprint(iface)
	}
	return defaultVal
}

// Slce returns slice of interface{} by key
// If key isn't defined or have a different type will be returned defaultVal arg value
func (s Map) Slice(key string, defaultVal []interface{}) (res []interface{}) {
	if arr, check := s[key].([]interface{}); check {
		res = arr
	} else {
		res = defaultVal
	}
	return
}

// StringSlice returns string slice by key
// If key isn't defined will be returned defaultVal arg value
func (s Map) StringSlice(key string, defaultVal []string) (res []string) {
	if arr, check := s[key].([]interface{}); check {
		res = make([]string, len(arr))
		for i, v := range arr {
			res[i] = fmt.Sprint(v)
		}
	} else {
		res = defaultVal
	}
	return
}

// Map returns Map object by key
// If key isn't defined or have other type will be returned defaultVal arg value
func (s Map) Map(key string, defaultVal Map) (res Map) {
	val := s[key]
	switch iface := val.(type) {
	case Map:
		res = iface
	case map[string]interface{}:
		res = Map(iface)
	default:
		res = defaultVal
	}
	return
}

// JSON Return JSON source of the self object
func (s Map) JSON() (res []byte) {
	res, _ = json.Marshal(s)
	return
}

// ToMap returns map[string]interface{} of the self object
func (s Map) ToMap() map[string]interface{} { return map[string]interface{}(s) }

// Copy returns copied map (todo dipcopy)
func (s Map) Copy() (res Map) {
	res = make(Map)
	for key, val := range s {
		res[key] = val
	}
	return
}

func (s Map) IsEmpty() bool {
	return len(s) == 0
}

func (s Map) Update(m Map) {
	for key, val := range m {
		s[key] = val
	}
}
