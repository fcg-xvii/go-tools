package concurrent

import (
	"sync"
)

type Element struct {
	m    *sync.RWMutex
	next *Element
	prev *Element
	val  interface{}
}

/*
func (s *Element) setPrev(elem *Element) {
	s.m.Lock()
	s.prev = elem
	s.m.Unlock()
}

func (s *Element) setNext(elem *Element) {
	s.m.Lock()
	s.next = elem
	s.m.Unlock()
}
*/

func (s *Element) SetVal(val interface{}) {
	s.m.Lock()
	s.val = val
	s.m.Unlock()
}

func (s *Element) Prev() (elem *Element) {
	s.m.RLock()
	elem = s.prev
	s.m.RUnlock()
	return
}

func (s *Element) Next() (elem *Element) {
	s.m.RLock()
	elem = s.next
	s.m.RUnlock()
	return
}

func (s *Element) Val() (val interface{}) {
	s.m.RLock()
	val = s.val
	s.m.RUnlock()
	return
}

func (s *Element) destroy() {
	s.m.Lock()
	s.prev, s.next = nil, nil
	s.m.Unlock()
}

func (s *Element) setNext(elem *Element) {
	s.m.Lock()
	if s.next != nil {
		elem.next, s.next.prev = s.next, elem
	}
	s.next, elem.prev = elem, s
	s.m.Unlock()
}

func (s *Element) setPref(elem *Element) {
	s.m.Lock()
	if s.prev != nil {

	}
	s.m.Unlock()
}

type List struct {
	m     *sync.RWMutex
	first *Element
	size  int
}

/*
func (s *List) insertNext(elem, base *Element) {
	s.m.Lock()
	if base == nil {
		// insert first element
		if s.first != nil {
			s.first.setNext(elem)
		}
		s.first = elem
	} else {
		base.setNext(elem)
	}
	s.m.Unlock()
}
*/

func (s *List) PushBack(v interface{}) *Element {
	elem := Element{
		val: v,
	}
	s.m.Lock()
	if s.first == nil {
		s.first = elem
	} else {

	}
	s.m.Lock()
}
