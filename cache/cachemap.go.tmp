package cache

import (
	"runtime"
	"sync"
	"time"
)

type CallCreate func(key interface{}) (value interface{}, created bool)
type CallCreateNew func(key interface{}) (rKey, value interface{}, created bool)
type CallCheck func(key, value interface{}, exists bool) (rKey, rValue interface{}, created bool)

func NewMap(liveDuration time.Duration, maxSize int) *CacheMap {
	res := &CacheMap{&cacheMap{
		locker:          new(sync.RWMutex),
		items:           make(map[interface{}]*cacheMapItem),
		liveDuration:    liveDuration,
		maxSize:         maxSize,
		stopCleanerChan: make(chan byte),
	}}
	runtime.SetFinalizer(res, destroyCacheMap)
	return res
}

type CacheMap struct {
	*cacheMap
}

type cacheMapItem struct {
	value  interface{}
	expire int64
}

type cacheMap struct {
	locker           *sync.RWMutex
	items            map[interface{}]*cacheMapItem
	liveDuration     time.Duration
	maxSize          int
	cleanerWork      bool
	stopCleanerChan  chan byte
	cleanedEventChan chan map[interface{}]interface{}
}

func (s *cacheMap) CleanEvent() (eventChan <-chan map[interface{}]interface{}) {
	s.locker.Lock()
	if s.cleanedEventChan == nil {
		s.cleanedEventChan = make(chan map[interface{}]interface{})
	}
	eventChan = s.cleanedEventChan
	s.locker.Unlock()
	return
}

func (s *cacheMap) runCleaner() {
	ticker := time.NewTicker(s.liveDuration / 2)
	for {
		select {
		case <-ticker.C:
			{
				now := time.Now().UnixNano()
				s.locker.Lock()
				cleaned := make(map[interface{}]interface{})
				for key, val := range s.items {
					if now > val.expire {
						cleaned[key] = val.value
						delete(s.items, key)
					}
				}
				if len(cleaned) > 0 {
					s.cleanedEventChan <- cleaned
				}
				if len(s.items) == 0 {
					s.cleanerWork = false
					ticker.Stop()
					s.locker.Unlock()
					return
				}
				s.locker.Unlock()
			}
		case <-s.stopCleanerChan:
			{
				s.cleanerWork = false
				ticker.Stop()
				return
			}
		}
	}
}

func (s *cacheMap) delete(key interface{}) {
	delete(s.items, key)
	if len(s.items) == 0 && s.cleanerWork {
		s.stopCleanerChan <- 0
	}
}

// Delete removes cached object
func (s *cacheMap) Delete(key interface{}) {
	s.locker.Lock()
	s.delete(key)
	s.locker.Unlock()
}

func (s *cacheMap) set(key, value interface{}) {
	s.items[key] = &cacheMapItem{
		value:  value,
		expire: time.Now().Add(s.liveDuration).UnixNano(),
	}
	if !s.cleanerWork {
		s.cleanerWork = true
		go s.runCleaner()
	}
}

func (s *cacheMap) Set(key, value interface{}) {
	s.locker.Lock()
	s.set(key, value)
	s.locker.Unlock()
}

func (s *cacheMap) SetMulti(m map[interface{}]interface{}) {
	s.locker.Lock()
	for key, val := range m {
		s.set(key, val)
	}
	s.locker.Unlock()
}

func (s *CacheMap) DeleteMulti(keys []interface{}) {
	s.locker.Lock()
	for _, key := range keys {
		delete(s.items, key)
	}
	s.locker.Unlock()
}

func (s *cacheMap) get(key interface{}) (res interface{}, check bool) {
	var item *cacheMapItem
	if item, check = s.items[key]; check {
		res = item.value
	}
	return
}

func (s *cacheMap) Get(key interface{}) (res interface{}, check bool) {
	s.locker.RLock()
	res, check = s.get(key)
	s.locker.RUnlock()
	return
}

func (s *cacheMap) GetCreate(key interface{}, mCreate CallCreate) (val interface{}, check bool) {
	if val, check = s.Get(key); !check {
		s.locker.Lock()
		if val, check = s.get(key); check {
			s.locker.Unlock()
			return
		}

		if val, check = mCreate(key); check {
			s.set(key, val)
		}
		s.locker.Unlock()
	}
	return
}

func (s *cacheMap) GetCreateNew(key interface{}, mCreateNew CallCreateNew) (rKey, val interface{}, check bool) {
	rKey = key
	if val, check = s.Get(key); !check {
		s.locker.Lock()
		if val, check = s.get(key); check {
			s.locker.Unlock()
			return
		}

		if rKey, val, check = mCreateNew(key); check {
			s.set(rKey, val)
		}
		s.locker.Unlock()
	}
	return
}

func (s *cacheMap) GetCheck(key interface{}, mCheck CallCheck) (res interface{}, check bool) {
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
func (s *cacheMap) Each(callback func(interface{}, interface{}) bool) {
	s.locker.RLock()
	for key, val := range s.items {
		if !callback(key, val.value) {
			s.locker.RUnlock()
			return
		}
	}
	s.locker.RUnlock()
}

func (s *cacheMap) Len() (res int) {
	s.locker.RLock()
	res = len(s.items)
	s.locker.RUnlock()
	return
}

func (s *cacheMap) Keys() (res []interface{}) {
	s.locker.RLock()
	res = make([]interface{}, 0, len(s.items))
	for key := range s.items {
		res = append(res, key)
	}
	s.locker.RUnlock()
	return
}

func (s *cacheMap) Clear() {
	s.locker.Lock()
	s.items = make(map[interface{}]*cacheMapItem)
	s.locker.Unlock()
}

// for garbage collector
func destroyCacheMap(m *CacheMap) {
	close(m.stopCleanerChan)
	close(m.cleanedEventChan)
}
