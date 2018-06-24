package memoryM

import (
	"container/heap"
	"time"
)

type advGC struct {
	h *heapStruct
}

func newGc() advGC {
	h := heapStruct{make([]*heapEl, 0), make(map[string]int)}
	heap.Init(&h)
	gc := advGC{h: &h}
	return gc
}

func (gc *advGC) set(key string, ttl int64) {
	if pos, ok := gc.h.dict[key]; ok {
		gc.h.store[pos].ttl = ttl
		heap.Fix(gc.h, pos)
		return
	}
	he := heapEl{key, ttl}
	heap.Push(gc.h, &he)
}

func (gc *advGC) del(key string) {
	pos, ok := gc.h.dict[key]
	delete(gc.h.dict, key)
	if !ok {
		return
	}
	heap.Remove(gc.h, pos)
}

func (gc *advGC) popExpired() []string {
	res := make([]string, 0)
	if len(gc.h.store) == 0 {
		return res
	}
	now := time.Now().Unix()
	for {
		he := gc.h.store[0]
		if now > he.ttl {
			res = append(res, he.key)
			heap.Pop(gc.h)
			delete(gc.h.dict, he.key)
		} else {
			return res
		}
	}
}

func (gc *advGC) hasExpired() bool {
	if len(gc.h.store) == 0 {
		return false
	}
	return time.Now().Unix() > gc.h.store[0].ttl
}