package memoryM

import (
	"errors"
	"sync"
	"time"
)

var ticker = time.Tick(time.Minute)

type Manager struct {
	mu    sync.RWMutex
	store map[string]el
	gs    advGC // garbage store
}

type el struct {
	val interface{}
	ttl int64
}

func New() *Manager {
	manager := Manager{
		sync.RWMutex{},
		make(map[string]el),
		newGc(),
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
	el := el{val, ttl.Unix()}
	man.mu.Lock()
	man.store[key] = el
	man.gs.set(key, el.ttl)
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
		return nil,  errors.New("not found")
	}
	return el.val, nil
}

func (man *Manager) Del(key string) error {
	if man.store == nil {
		return errors.New("manager is not initialised")
	}
	man.mu.Lock()
	delete(man.store, key)
	man.gs.del(key)
	man.mu.Unlock()
	return nil
}


func (man *Manager) simpleGC() {
	now := time.Now().Unix()
	man.mu.RLock() // locks writers
	exp := make([]string, 0)
	for key, el := range man.store {
		if now > el.ttl {
			exp = append(exp, key)
		}
	}
	man.mu.RUnlock()
	if len(exp) == 0 {
		return
	}
	// something may change here, so we need re-check expiration
	man.mu.Lock() // locks writers + readers
	for _, key := range exp {
		el, ok := man.store[key]
		if ok && now > el.ttl {
			delete(man.store, key)
		}
	}
	man.mu.Unlock()
}

func (man *Manager) advGC() {
	man.mu.RLock()
	if ! man.gs.hasExpired() {
		man.mu.RUnlock()
		return
	}
	man.mu.RUnlock()
	man.mu.Lock()
	exp := man.gs.popExpired()
	for _, toDel := range exp {
		delete(man.store, toDel)
	}
	man.mu.Unlock()
}
