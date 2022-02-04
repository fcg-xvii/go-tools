package concurrent

import (
	"sync"
)

func initElem(v interface{}) *Element {
	return &Element{
		m:   new(sync.RWMutex),
		val: v,
	}
}

type Element struct {
	m    *sync.RWMutex
	prev *Element
	next *Element
	val  interface{}
}

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
	if s.prev != nil {
		s.prev.next = s.next
	}
	if s.next != nil {
		s.next.prev = s.prev
	}
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

func (s *Element) setPrev(elem *Element) {
	s.m.Lock()
	if s.prev != nil {
		elem.prev, s.prev.next = s.prev, elem
	}
	s.prev, elem.next = elem, s
	s.m.Unlock()
}

func ListFromSlise(sl []interface{}) *List {
	l := NewList()
	for _, val := range sl {
		l.PushBack(val)
	}
	return l
}

func NewList() *List {
	return &List{
		m: new(sync.RWMutex),
	}
}

type List struct {
	m     *sync.RWMutex
	first *Element
	last  *Element
	size  int
}

func (s *List) PushBackIfNotExists(v interface{}) (elem *Element, added bool) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.search(v) == nil {
		added = true
		elem = s.pushBack(v)
	}
	return
}

func (s *List) PushBack(v interface{}) (elem *Element) {
	s.m.Lock()
	elem = s.pushBack(v)
	s.m.Unlock()
	return
}

func (s *List) pushBack(v interface{}) *Element {
	elem := initElem(v)
	if s.first == nil {
		s.first, s.last = elem, elem
	} else {
		s.last.setNext(elem)
		s.last = elem
	}
	s.size++
	return elem
}

func (s *List) PushFront(v interface{}) (elem *Element) {
	s.m.Lock()
	elem = s.pushFront(v)
	s.m.Unlock()
	return
}

func (s *List) PushFrontIfNotExists(v interface{}) (elem *Element, added bool) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.search(v) == nil {
		added = true
		elem = s.pushFront(v)
	}
	return
}

func (s *List) pushFront(v interface{}) *Element {
	elem := initElem(v)
	if s.first == nil {
		s.first, s.last = elem, elem
	} else {
		s.first.setPrev(elem)
		s.first = elem
	}
	s.size++
	return elem
}

func (s *List) Remove(elem *Element) {
	if elem == nil {
		return
	}
	s.m.Lock()
	if elem == s.first {
		s.first = elem.next
	}
	if elem == s.last {
		s.last = elem.prev
	}
	elem.destroy()
	s.size--
	s.m.Unlock()
	return
}

func (s *List) Size() (size int) {
	s.m.RLock()
	size = s.size
	s.m.RUnlock()
	return
}

func (s *List) First() (elem *Element) {
	s.m.RLock()
	elem = s.first
	s.m.RUnlock()
	return
}

func (s *List) Last() (elem *Element) {
	s.m.RLock()
	elem = s.last
	s.m.RUnlock()
	return
}

func (s *List) Index(index int) (elem *Element) {
	i := 0
	elem = s.First()
	for elem != nil && i < index {
		elem, i = elem.Next(), i+1
	}
	return
}

func (s *List) Slice() []interface{} {
	s.m.RLock()
	f, res := s.first, make([]interface{}, 0, s.size)
	s.m.RUnlock()
	for f != nil {
		res, f = append(res, f.Val()), f.Next()
	}
	return res
}

func (s *List) search(val interface{}) *Element {
	for f := s.first; f != nil; f = f.Next() {
		if f.Val() == val {
			return f
		}
	}
	return nil
}

func (s *List) Search(val interface{}) *Element {
	for f := s.First(); f != nil; f = f.Next() {
		if f.Val() == val {
			return f
		}
	}
	return nil
}
