package concurrent

import (
	"sync"

	"github.com/fcg-xvii/go-tools/containers"
)

func NewQueue() *Queue {
	return &Queue{
		q: containers.NewQueue(),
		m: new(sync.RWMutex),
	}
}

type Queue struct {
	q *containers.Queue
	m *sync.RWMutex
}

func (s *Queue) Size() (size int) {
	s.m.RLock()
	size = s.q.Size()
	s.m.RUnlock()
	return
}

func (s *Queue) Push(val interface{}) {
	s.m.Lock()
	s.q.Push(val)
	s.m.Unlock()
}

func (s *Queue) Pop() (val interface{}, check bool) {
	s.m.Lock()
	val, check = s.q.Pop()
	s.m.Unlock()
	return
}
