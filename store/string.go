package store

import "sync"

type StringCallCreate func(key string) (value interface{}, created bool)
type StringCallCreateMulti func(key string) (map[string]interface{}, bool)
type StringCallCheck func(key string, value interface{}, exists bool) (rKey string, rValue interface{}, created bool)

func StringFromMap(m map[string]interface{}) *StoreString {
	return &StoreString{
		locker: new(sync.RWMutex),
		items:  m,
	}
}

func StringNew() *StoreString {
	return &StoreString{
		locker: new(sync.RWMutex),
		items:  make(map[string]interface{}),
	}
}

type StoreString struct {
	locker *sync.RWMutex
	items  map[string]interface{}
}

func (s *StoreString) delete(key string) {
	delete(s.items, key)
}

func (s *StoreString) Delete(key string) {
	s.locker.Lock()
	delete(s.items, key)
	s.locker.Unlock()
}

func (s *StoreString) DeleteMulti(keys []string) {
	s.locker.Lock()
	for _, key := range keys {
		delete(s.items, key)
	}
	s.locker.Unlock()
}

func (s *StoreString) set(key string, val interface{}) {
	s.items[key] = val
}

func (s *StoreString) setMulti(m map[string]interface{}) {
	for key, val := range m {
		s.set(key, val)
	}
}

func (s *StoreString) Set(key string, val interface{}) {
	s.locker.Lock()
	s.set(key, val)
	s.locker.Unlock()
}

func (s *StoreString) SetMulti(m map[string]interface{}) {
	s.locker.Lock()
	s.setMulti(m)
	s.locker.Unlock()
}

func (s *StoreString) get(key string) (val interface{}, check bool) {
	val, check = s.items[key]
	return
}

func (s *StoreString) Get(key string) (val interface{}, check bool) {
	s.locker.RLock()
	val, check = s.get(key)
	s.locker.RUnlock()
	return
}

func (s *StoreString) GetCreate(key string, mCreate StringCallCreate) (res interface{}, check bool) {
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

func (s *StoreString) GetCreateMulti(key string, mCreateMulti StringCallCreateMulti) (res interface{}, check bool) {
	if res, check = s.Get(key); !check {
		s.locker.Lock()
		if res, check = s.get(key); check {
			s.locker.Unlock()
			return
		}

		var m map[string]interface{}
		if m, check = mCreateMulti(key); check {
			s.setMulti(m)
			res, check = s.items[key]
		}
		s.locker.Unlock()
	}
	return
}

func (s *StoreString) GetCheck(key string, mCheck StringCallCheck) (res interface{}, check bool) {
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
func (s *StoreString) Each(callback func(string, interface{}) bool) {
	s.locker.RLock()
	for key, val := range s.items {
		if !callback(key, val) {
			s.locker.RUnlock()
			return
		}
	}
	s.locker.RUnlock()
}

func (s *StoreString) Len() (res int) {
	s.locker.RLock()
	res = len(s.items)
	s.locker.RUnlock()
	return
}

func (s *StoreString) Keys() (res []string) {
	s.locker.RLock()
	res = make([]string, 0, len(s.items))
	for key := range s.items {
		res = append(res, key)
	}
	s.locker.RUnlock()
	return
}

func (s *StoreString) Clear() {
	s.locker.Lock()
	s.items = make(map[string]interface{})
	s.locker.Unlock()
}

func (s *StoreString) Map() (res map[string]interface{}) {
	res = make(map[string]interface{})
	s.locker.RLock()
	for key, val := range s.items {
		res[key] = val
	}
	s.locker.RUnlock()
	return
}
