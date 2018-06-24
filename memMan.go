package memMan

import (
	"container/heap"
	"errors"
	"sync"
	"time"
	"math"
)

var ticker = time.Tick(time.Minute)

type Manager struct {
	mu    sync.RWMutex
	store map[string]*element
	hs    *heapStruct
}

type element struct {
	key string
	val interface{}
	ttl int64
	pos int
}

func New() *Manager {
	manager := Manager{
		sync.RWMutex{},
		make(map[string]*element),
		&heapStruct{make([]*element, 0)},
	}
	go func() {
		for {
			<-ticker
			manager.advGC()
		}
	}()
	return &manager
}

func (man *Manager) Set(key string, val interface{}, ttl time.Time) error {
	if man.store == nil {
		return errors.New("manager is not initialised")
	}
	if len(man.hs.store) == math.MaxInt64 {
		return errors.New("manager is full")
	}
	man.mu.Lock()
	el, exists := man.store[key]
	if exists {
		el.val = val
		el.ttl = ttl.Unix()
		heap.Fix(man.hs, el.pos)
		man.mu.Unlock()
		return nil
	}
	el = &element{key, val, ttl.Unix(), 0}
	man.store[key] = el
	heap.Push(man.hs, el)
	man.mu.Unlock()
	return nil
}

func (man *Manager) Get(key string) (interface{}, error) {
	if man.store == nil {
		return nil, errors.New("manager is not initialised")
	}
	man.mu.RLock()
	el, ok := man.store[key]
	man.mu.RUnlock()
	if !ok || time.Now().Unix() > el.ttl {
		return nil, errors.New("not found")
	}
	return el.val, nil
}

func (man *Manager) Del(key string) error {
	if man.store == nil {
		return errors.New("manager is not initialised")
	}
	man.mu.Lock()
	elem, exists := man.store[key]
	if exists {
		heap.Remove(man.hs, elem.pos)
		delete(man.store, key)
	}
	man.mu.Unlock()
	return nil
}

func (man *Manager) superSimpleGC() {
	if len(man.store) == 0 {
		return
	}
	now := time.Now().Unix()
	man.mu.Lock() // locks writers + readers
	for key, el := range man.store {
		if now > el.ttl {
			delete(man.store, key)
		}
	}
	man.mu.Unlock()
}

func (man *Manager) simpleGC() {
	if len(man.store) == 0 {
		return
	}
	now := time.Now().Unix()
	man.mu.RLock() // locks writers
	expKeys := make([]string, 0)
	for key, el := range man.store {
		if now > el.ttl {
			expKeys = append(expKeys, key)
		}
	}
	man.mu.RUnlock()
	if len(expKeys) == 0 {
		return
	}
	// something may change here, so we need re-check expiration
	man.mu.Lock() // locks writers + readers
	for _, key := range expKeys {
		el, ok := man.store[key]
		if ok && now > el.ttl {
			delete(man.store, key)
		}
	}
	man.mu.Unlock()
}

func (man *Manager) advGC() {
	man.mu.RLock()
	if len(man.store) == 0 {
		return
	}
	now := time.Now().Unix()
	topEl := man.hs.store[0]
	man.mu.RUnlock()
	if now < topEl.ttl {
		return
	}
	man.mu.Lock()
	for {
		if len(man.store) == 0 {
			man.mu.Unlock()
			return
		}
		topEl = man.hs.store[0]
		if now > topEl.ttl {
			delete(man.store, topEl.key)
			heap.Pop(man.hs)
		} else {
			man.mu.Unlock()
			return
		}
	}
}
