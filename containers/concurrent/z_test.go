package concurrent

import (
	"log"
	"testing"
)

func TestQueue(t *testing.T) {
	q := NewQueue()
	q.Push("one", "two", "three")
	q.Push("four")
	q.Push("five")
	t.Log(q)
	for {
		if val, check := q.Pop(); check {
			t.Log(val, q.Size())
		} else {
			break
		}
	}
	t.Log(q)
}

func showList(l *List) {
}

func TestList(t *testing.T) {
	l := NewList()
	log.Println(l)
	for i := 0; i < 10; i++ {
		l.PushFront(i)
	}
	log.Println(l.Slice(), l.Size())
	elem := l.Index(0)
	log.Println(elem)
	l.Remove(elem)
	log.Println(l.Slice())
	for l.Size() > 0 {
		l.Remove(l.First())
		log.Println(l.Slice(), l.Size())
	}
	log.Println(l.first, l.last)
	l.PushBack(100)
	log.Println(l.Slice(), l.Size(), l.first, l.last)

	sl := make([]interface{}, 10)
	for i := range sl {
		sl[i] = i * 100
	}
	log.Println(sl)
	l = ListFromSlise(sl)
	log.Println(l.Slice())

	elem = l.Search(500)
	log.Println(elem)
	elem = l.Search(1500)
	log.Println(1500, "search", elem)

	elem, added := l.PushBackIfNotExists(500)
	log.Println(500, elem, added)
	elem, added = l.PushBackIfNotExists(1500)
	log.Println(1500, elem, added)

}
