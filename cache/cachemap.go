package cache

import (
	"sync"
	"time"
)

type cacheMapItem struct {
	value  interface{}
	expire int64
}

type cacheMap struct {
	locker          *sync.RWMutex
	items           map[interface{}]*cacheMapItem
	liveDuration    time.Duration
	maxSize         int64
	cleanerWork     bool
	stopCleanerChan chan struct{}
}

func (s *cacheMap) runCleaner() {
	ticker := time.NewTicker(s.liveDuration / 2)
	for {
		select {
		case <-ticker.C:
			{
				now := time.Now().UnixNano()
				s.locker.Lock()
				for key, val := range s.items {
					if now > val.expire {
						s.delete(key)
					}
				}
				s.locker.RUnlock()
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

type CreateCall func(key interface{}) (interface{}, bool)

func (s *cacheMap) get(key interface{}, call CreateCall) (res interface{}, check bool) {
	var item *cacheMapItem
	if item, check = s.items[key]; !check && call != nil {
		if cVal, cCheck := call(key); cCheck {
			s.set(key, cVal)
		}
	} else {
		///...
	}
	return
}

func (s *cacheMap) Get(key interface{}) (res interface{}, check bool) {
	s.locker.RLock()
	res, check = s.get(key, cCall)
	s.locker.RUnlock()
	return
}

func (s *cacheMap) GetOrCreate(key interface{}, CreateMethod) interface{}) (res interface{}, check bool) {
	if res, check = s.Get(key, cCall); !check {
		s.locker.Lock()
		if res, check = s.get(key, cCall); check {
			s.locker.Unlock()
			return
		}

		if key, res, check = createCall(key); check {
			s.set(key, res)
		}
		s.locker.Unlock()
	}
	return
}

// Delete removes cached object
func (s *cacheMap) Delete(key interface{}) {
	s.locker.Lock()
	s.delete(key)
	s.locker.Unlock()
}

// Size returns cache map size
func (s *cacheMap) Size() (res int) {
	s.locker.RLock()
	res = len(s.items)
	s.locker.RUnlock()
	return res
}

// Each implements a map bypass for each key using the callback function. If the callback function returns false, then the cycle stops
func (s *cacheMap) Each(callback func(interface{}, interface{}) bool) {
	s.locker.RLock()
	for key, val := range s.items {
		if !callback(key, val) {
			s.locker.Unlock()
			return
		}
	}
	s.locker.Unlock()
}

// for garbage collector
func destroyCacheMap(m *cacheMap) {
	close(m.stopCleanerChan)
}
