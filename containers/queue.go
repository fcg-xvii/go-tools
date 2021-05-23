package containers

type queueItem struct {
	val  interface{}
	next *queueItem
}

func NewQueue() *Queue {
	return new(Queue)
}

type Queue struct {
	start *queueItem
	top   *queueItem
	size  int
}

func (s *Queue) Size() int {
	return s.size
}

func (s *Queue) Push(val interface{}) {
	item := &queueItem{
		val: val,
	}
	if s.top != nil {
		s.top.next = item
	} else {
		s.start = item
	}
	s.top = item
	s.size++
}

func (s *Queue) Pop() (val interface{}, check bool) {
	if s.start == nil {
		return
	}
	val, check = s.start.val, true
	s.start = s.start.next
	if s.start == nil {
		s.top = nil
	}
	s.size--
	return
}
