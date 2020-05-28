// Package jsonmap for quickly get typed value from map[string]interface{} object
package jsonmap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// JSONMap type
type JSONMap map[string]interface{}

// New init JSONMap object
func New() JSONMap {
	return make(JSONMap)
}

// FromMap convert map[string]interface{} to JSONMap object
func FromMap(m map[string]interface{}) JSONMap {
	return JSONMap(m)
}

// Bool returns bool value by key
func (s JSONMap) Bool(key string, defaultVal bool) bool {
	if res, check := s[key].(bool); check {
		return res
	}
	return defaultVal
}

// KeyExists check value exists by key
func (s JSONMap) KeyExists(key string) bool {
	_, check := s[key]
	return check
}

// KeysExists check values exists by keys list.
// Returns the first blank key found.
// If all keys are defined, an empty string will be returned
func (s JSONMap) KeysExists(keys []string) string {
	for _, key := range keys {
		if _, check := s[key]; !check {
			return key
		}
	}
	return ""
}

// Int returns int64 value by key.
// If key isn't defined will be returned defaultVal arg value
func (s JSONMap) Int(key string, defaultVal int64) int64 {
	if iface, check := s[key]; check {
		rVal := reflect.ValueOf(iface)
		switch rVal.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return rVal.Int()
		case reflect.Float32, reflect.Float64:
			return int64(rVal.Float())
		}
	}
	return defaultVal
}

// Value returns interface object with attempt to convert to defaultVal type.
// If key isn't defined will be returned defaultVal arg value
func (s JSONMap) Value(key string, defaultVal interface{}) interface{} {

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
func (s JSONMap) ValueJSON(key string, defaultVal []byte) (res []byte) {
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
func (s JSONMap) String(key, defaultVal string) string {
	if iface, check := s[key]; check {
		return fmt.Sprint(iface)
	}
	return defaultVal
}

// Slce returns slice of interface{} by key
// If key isn't defined or have a different type will be returned defaultVal arg value
func (s JSONMap) Slice(key string, defaultVal []interface{}) (res []interface{}) {
	if arr, check := s[key].([]interface{}); check {
		res = arr
	} else {
		res = defaultVal
	}
	return
}

// StringSlice returns string slice by key
// If key isn't defined will be returned defaultVal arg value
func (s JSONMap) StringSlice(key string, defaultVal []string) (res []string) {
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

// JSONMap returns JSONMap object by key
// If key isn't defined or have other type will be returned defaultVal arg value
func (s JSONMap) JSONMap(key string, defaultVal JSONMap) (res JSONMap) {
	if m, check := s[key].(map[string]interface{}); check {
		res = JSONMap(m)
	} else {
		res = defaultVal
	}
	return
}

// JSON Return JSON source of the self object
func (s JSONMap) JSON() (res []byte) {
	res, _ = json.Marshal(s)
	return
}

// ToMap returns map[string]interface{} of the self object
func (s JSONMap) ToMap() map[string]interface{} { return map[string]interface{}(s) }
