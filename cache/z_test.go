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
}
