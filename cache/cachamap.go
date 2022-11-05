package cache

import (
	"context"
	"sync"
	"time"
)

type item struct {
	value  any
	expire int64
}

type CacheMap struct {
	finished   bool
	ctx        context.Context
	items      map[any]*item
	mi         *sync.RWMutex
	live       time.Duration
	maxSize    int
	cleanChans []chan map[any]any
	mc         *sync.RWMutex
}

func (s *CacheMap) IsFinished() bool {
	return s.finished
}

func (s *CacheMap) EventCleaner() (ch <-chan map[any]any) {
	s.mc.Lock()
	ech := make(chan map[any]any, 1)
	s.cleanChans = append(s.cleanChans, ech)
	ch = ech
	s.mc.Unlock()
	return
}
