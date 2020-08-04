package cache

import (
	"fmt"
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
		for i := 0; i < 1000; i++ {
			go func(j int) {
				m.GetOrCreate(fmt.Sprintf("key%v", j), func(interface{}) (interface{}, bool) {
					//log.Println("CREATE_CALL")
					return 100, true
				})
			}(i)
			time.Sleep(time.Millisecond * 5)
			//log.Println("i", i)
		}

	}()

	time.Sleep(time.Second * 3)
	log.Println("!!!")
	for {
		log.Println("len", m.Len())
		/*if m.Len() == 0 {
			break
		}*/
		time.Sleep(time.Millisecond * 500)
	}

	log.Println(m.Get("key1"))
}
