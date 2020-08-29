package store

import (
	"log"
	"testing"
)

func TestStore(t *testing.T) {
	st := New()

	st.Set("one", 10)

	log.Println(st.Get("one"))
	log.Println(st.Get("two"))

	log.Println(st.GetCreate("two", func(key interface{}) (value interface{}, created bool) {
		return 20, true
	}))

	st.Delete("two")

	log.Println(st.GetCreate("two", func(key interface{}) (value interface{}, created bool) {
		return 30, true
	}))

	log.Println(st.Map())

	st = FromMap(map[interface{}]interface{}{
		"one": 1,
		"two": 2,
	})

	log.Println(st.Map())
}

func TestStoreString(t *testing.T) {
	st := StringNew()

	st.Set("one", 10)

	log.Println(st.Get("one"))
	log.Println(st.Get("two"))

	log.Println(st.GetCreate("two", func(key string) (value interface{}, created bool) {
		return 20, true
	}))

	st.Delete("two")

	log.Println(st.GetCreate("two", func(key string) (value interface{}, created bool) {
		return 30, true
	}))

	log.Println(st.Map())

	st = StringFromMap(map[string]interface{}{
		"one": 1,
		"two": 2,
	})

	log.Println(st.Map())

	val, check := st.GetCreateMulti("key1", func(string) (map[string]interface{}, bool) {
		return map[string]interface{}{
			"key100": 100,
			"key2":   200,
		}, true
	})

	log.Println(val, check)
	log.Println("==================")
	log.Println(st.items)
}
