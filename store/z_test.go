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
}
