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
	locker               *sync.RWMutex
	items                map[interface{}]*cacheMapItem
	liveDuration         time.Duration
	maxSize              int64
	cleanerWork          bool
	stopCleanerChan      chan struct{}
	itemRewmovedCallback func(key, value interface{})
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
						s.delete(key, val)
					}
				}
				s.locker.RUnlock()
			}
		case <-s.stopCleanerChan:
			{
				s.cleanerWork = false
				ticker.Stop()
				// todo...
				/*if s.callbackRemoved != nil {
					s.callbackRemoved(s.items)
				}*/
				return
			}
		}
	}
}

func (s *cacheMap) delete(key interface{}, value ...interface{}) {
	if s.itemRewmovedCallback != nil {
		var val interface{}
		if len(value) > 0 {
			val = value[0]
		} else {
			val = s.items[key]
		}
		s.itemRewmovedCallback(key, val)
	}
	delete(s.items, key)
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
