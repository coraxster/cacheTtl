package cacheTtl

import (
	"container/heap"
	"errors"
	"math"
	"sync"
	"time"
)

type Cache struct {
	mu    sync.RWMutex
	store map[string]*element
	hs    *heapStruct
}

type element struct {
	key   string
	val   interface{}
	ttl   int64
	index int // position in heap
}

var ErrNotFound = errors.New("not found")

func New() *Cache {
	hs := heapStruct(make([]*element, 0))
	c := Cache{
		sync.RWMutex{},
		make(map[string]*element),
		&hs,
	}
	ticker := time.Tick(time.Minute)
	go func() {
		for {
			<-ticker
			c.advGC()
		}
	}()
	return &c
}

func (c *Cache) Set(key string, val interface{}, ttl time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(*c.hs) == math.MaxInt64 {
		return errors.New("cache is full")
	}
	el, exists := c.store[key]
	if exists {
		el.val = val
		if el.ttl != ttl.Unix() {
			el.ttl = ttl.Unix()
			heap.Fix(c.hs, el.index)
		}
		return nil
	}
	el = &element{key, val, ttl.Unix(), 0}
	c.store[key] = el
	heap.Push(c.hs, el)
	return nil
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	el, ok := c.store[key]
	if !ok || time.Now().Unix() > el.ttl {
		return nil, ErrNotFound
	}
	return el.val, nil
}

func (c *Cache) Del(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, exists := c.store[key]
	if exists {
		heap.Remove(c.hs, elem.index)
		delete(c.store, key)
	}
	return nil
}

func (c *Cache) superSimpleGC() {
	c.mu.Lock() // locks writers + readers
	defer c.mu.Unlock()
	if len(c.store) == 0 {
		return
	}
	now := time.Now().Unix()
	for key, el := range c.store {
		if now > el.ttl {
			delete(c.store, key)
		}
	}
}

func (c *Cache) simpleGC() {
	c.mu.RLock() // locks writers
	if len(c.store) == 0 {
		return
	}
	now := time.Now().Unix()
	expKeys := make([]string, 0)
	for key, el := range c.store {
		if now > el.ttl {
			expKeys = append(expKeys, key)
		}
	}
	c.mu.RUnlock()
	if len(expKeys) == 0 {
		return
	}
	// something may change here, so we need re-check expiration
	c.mu.Lock() // locks writers + readers
	defer c.mu.Unlock()
	for _, key := range expKeys {
		el, ok := c.store[key]
		if ok && now > el.ttl {
			delete(c.store, key)
		}
	}
}

func (c *Cache) advGC() {
	c.mu.RLock()
	if len(c.store) == 0 {
		return
	}
	now := time.Now().Unix()
	topEl := (*c.hs)[0]
	if now < topEl.ttl {
		c.mu.RUnlock()
		return
	}
	c.mu.RUnlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	for {
		if len(c.store) == 0 {
			return
		}
		topEl = (*c.hs)[0]
		if now > topEl.ttl {
			delete(c.store, topEl.key)
			heap.Pop(c.hs)
		} else {
			return
		}
	}
}
