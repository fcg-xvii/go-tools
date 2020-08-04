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
