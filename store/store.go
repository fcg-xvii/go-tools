package store

import (
	"sync"
)

type CallCreate func(key interface{}) (value interface{}, created bool)
type CallCreateMulti func(key interface{}) (m map[interface{}]interface{}, created bool)
type CallCheck func(key, value interface{}, exists bool) (rKey, rValue interface{}, created bool)

func FromMap(m map[interface{}]interface{}) *Store {
	return &Store{
		locker: new(sync.RWMutex),
		items:  m,
	}
}

func New() *Store {
	return &Store{
		locker: new(sync.RWMutex),
		items:  make(map[interface{}]interface{}),
	}
}

type Store struct {
	locker *sync.RWMutex
	items  map[interface{}]interface{}
}

func (s *Store) delete(key interface{}) {
	delete(s.items, key)
}

func (s *Store) Delete(key interface{}) {
	s.locker.Lock()
	delete(s.items, key)
	s.locker.Unlock()
}

func (s *Store) DeleteMulti(keys []interface{}) {
	s.locker.Lock()
	for _, key := range keys {
		delete(s.items, key)
	}
	s.locker.Unlock()
}

func (s *Store) set(key, val interface{}) {
	s.items[key] = val
}

func (s *Store) setMulti(m map[interface{}]interface{}) {
	for key, val := range m {
		s.items[key] = val
	}
}

func (s *Store) Set(key, val interface{}) {
	s.locker.Lock()
	s.set(key, val)
	s.locker.Unlock()
}

func (s *Store) SetMulti(m map[interface{}]interface{}) {
	s.locker.Lock()
	s.setMulti(m)
	s.locker.Unlock()
}

func (s *Store) get(key interface{}) (val interface{}, check bool) {
	val, check = s.items[key]
	return
}

func (s *Store) Get(key interface{}) (val interface{}, check bool) {
	s.locker.RLock()
	val, check = s.get(key)
	s.locker.RUnlock()
	return
}

func (s *Store) GetCreate(key interface{}, mCreate CallCreate) (res interface{}, check bool) {
	if res, check = s.Get(key); !check {
		s.locker.Lock()
		if res, check = s.get(key); check {
			s.locker.Unlock()
			return
		}

		if res, check = mCreate(key); check {
			s.set(key, res)
		}
		s.locker.Unlock()
	}
	return
}

func (s *Store) GetCreateMulti(key interface{}, mCreateMulti CallCreateMulti) (res interface{}, check bool) {
	if res, check = s.Get(key); !check {
		s.locker.Lock()
		if res, check = s.get(key); check {
			s.locker.Unlock()
			return
		}

		var m map[interface{}]interface{}
		if m, check = mCreateMulti(key); check {
			s.setMulti(m)
			res, check = s.items[key]
		}
		s.locker.Unlock()
	}
	return
}

func (s *Store) GetCheck(key interface{}, mCheck CallCheck) (res interface{}, check bool) {
	s.locker.Lock()
	res, check = s.get(key)
	if rKey, rVal, rCheck := mCheck(key, res, check); rCheck {
		s.set(rKey, rVal)
		key, res, check = rKey, rVal, true
	}
	s.locker.Unlock()
	return
}

// Each implements a map bypass for each key using the callback function. If the callback function returns false, then the cycle stops
func (s *Store) Each(callback func(interface{}, interface{}) bool) {
	s.locker.RLock()
	for key, val := range s.items {
		if !callback(key, val) {
			s.locker.RUnlock()
			return
		}
	}
	s.locker.RUnlock()
}

func (s *Store) Len() (res int) {
	s.locker.RLock()
	res = len(s.items)
	s.locker.RUnlock()
	return
}

func (s *Store) Keys() (res []interface{}) {
	s.locker.RLock()
	res = make([]interface{}, 0, len(s.items))
	for key := range s.items {
		res = append(res, key)
	}
	s.locker.RUnlock()
	return
}

func (s *Store) Clear() {
	s.locker.Lock()
	s.items = make(map[interface{}]interface{})
	s.locker.Unlock()
}

func (s *Store) Map() (res map[interface{}]interface{}) {
	res = make(map[interface{}]interface{})
	s.locker.RLock()
	for key, val := range s.items {
		res[key] = val
	}
	s.locker.RUnlock()
	return
}
