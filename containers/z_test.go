package containers

import "testing"

func TestQueue(t *testing.T) {
	var q Queue
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
