package cache

import (
	"log"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	m := NewMap(time.Second*5, 10)

	m.Set("key", 10)

	log.Println(m.Get("key"))
	m.Delete("key")
	log.Println(m.Get("key"))
	log.Println(m.Get("key1"))

	go func() {
		for e := range m.CleanEvent() {
			log.Println("CLEANED", e)
		}
	}()

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second * 10)
		m.Set("key", "oneee")

	}
}

func TestMapCheck(t *testing.T) {
	m := NewMap(time.Hour, 0)

	m.Set("key", 10)

	log.Println(m.GetCheck("key", func(key, value interface{}, exists bool) (rKey, rVal interface{}, needUpdate bool) {
		log.Println("MCHECK", key, value, exists)
		rVal = "Cool!!!"
		rKey, needUpdate = "key1", true
		return
	}))

	log.Println("===================")
	log.Println(m.Get("key"))
	log.Println(m.Get("key1"))

	log.Println(m.GetCreateNew("key", func(key interface{}) (rKey, val interface{}, created bool) {
		return "oooone", 1000, true
	}))
}

func TestMapEach(t *testing.T) {
	m := NewMap(time.Second*5, 10)
	m.SetMulti(map[interface{}]interface{}{
		1: "one",
		2: "two",
		3: "three",
	})

	m.Each(func(key, val interface{}) bool {
		log.Println(key, val)
		return false
	})
}

func TestKeys(t *testing.T) {
	m := NewMap(time.Second*5, 20)
	m.SetMulti(map[interface{}]interface{}{
		"one": 1,
		"two": 2,
	})
	log.Println(m.Keys())
	m.Clear()
	log.Println(m.Keys())

}
