package jsonmap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// JSONMap
type JSONMap map[string]interface{}

// Getter of bool value
func (s JSONMap) Bool(key string, defaultVal bool) bool {
	if res, check := s[key].(bool); check {
		return res
	}
	return defaultVal
}

func (s JSONMap) KeyExists(key string) bool {
	_, check := s[key]
	return check
}

// Getter of int64 value
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

// Getter of value with attempt to convert to defaultVal type. If defaultVal is invalid (null) then will be return found value with each type
func (s JSONMap) Value(key string, defaultVal interface{}) interface{} {

	// check value exists by key
	if iface, check := s[key]; check {

		// check defaultVal is valid
		dVal := reflect.ValueOf(defaultVal)
		if !dVal.IsValid() {
			// invalid, return arrived interface
			return iface
		} else {
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
							if val, err := strconv.Atoi(lVal.String()); err != nil {
								return defaultVal
							} else {
								lVal = reflect.ValueOf(val)
								return lVal.Convert(dVal.Type())
							}
						case reflect.Float32, reflect.Float64:
							if val, err := strconv.ParseFloat(lVal.String(), 64); err != nil {
								return defaultVal
							} else {
								lVal = reflect.ValueOf(val)
								return lVal.Convert(dVal.Type())
							}
						}
					}
				}
			}
		}
	}
	return defaultVal
}

func (s JSONMap) String(key, defaultVal string) string {
	if iface, check := s[key]; check {
		return fmt.Sprint(iface)
	}
	return defaultVal
}

func (s JSONMap) StringArray(key string) (res []string) {
	if arr, check := s[key].([]interface{}); check {
		res = make([]string, len(arr))
		for i, v := range arr {
			res[i] = fmt.Sprint(v)
		}
	}
	return
}

func (s JSONMap) JSONMap(key string) (res JSONMap) {
	if m, check := s[key].(map[string]interface{}); check {
		res = JSONMap(m)
	}
	return
}

func (s JSONMap) JSON() (res []byte) {
	res, _ = json.Marshal(s)
	return
}

func (s JSONMap) ToMap() map[string]interface{} { return map[string]interface{}(s) }
