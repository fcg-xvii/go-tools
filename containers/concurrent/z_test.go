package concurrent

import (
	"log"
	"testing"
)

func TestQueue(t *testing.T) {
	q := NewQueue()
	q.Push("one")
	q.Push("two")
	q.Push("three")
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
}
