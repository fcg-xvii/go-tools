package store

import (
	"log"
	"reflect"
	"testing"
	"unsafe"
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

	var i int
	log.Printf("Size of var (reflect.TypeOf.Size): %d\n", reflect.TypeOf(i).Size())
	log.Printf("Size of var (unsafe.Sizeof): %d\n", unsafe.Sizeof(i))
}
