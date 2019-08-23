package cacheTtl

import (
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	k := "key"
	v := "value"
	cache := New()
	cache.Set(k, v, time.Now().Add(time.Minute))
	mv, er := cache.Get(k)
	if er != nil {
		t.Error(er)
	}
	if mv != v {
		t.Error("output is not eq input")
	}
}

func TestDel(t *testing.T) {
	k := "key"
	v := "value"
	cache := New()
	cache.Set(k, v, time.Now().Add(time.Minute))
	cache.Del(k)
	mv, er := cache.Get(k)
	if er != ErrNotFound {
		t.Error(er)
	}
	if mv != nil {
		t.Error("output is not nil after del")
	}
}

func TestTtl(t *testing.T) {
	k1 := "key1"
	v1 := "value"
	k2 := "key2"
	v2 := "value"
	cache := New()
	cache.Set(k1, v1, time.Now().Add(-time.Minute))
	cache.Set(k2, v2, time.Now().Add(time.Minute))
	mv1, er := cache.Get(k1)
	if er != ErrNotFound {
		t.Error(er)
	}
	if mv1 != nil {
		t.Error("ttl expired, output is not nil")
	}
	mv2, er := cache.Get(k2)
	if mv2 != v2 {
		t.Error("output is not eq input")
	}
}

func TestTtlByTime(t *testing.T) {
	k1 := "key1"
	v1 := "value"
	k2 := "key2"
	v2 := "value"
	cache := New()
	cache.Set(k1, v1, time.Now().Add(-time.Minute))
	cache.Set(k2, v2, time.Now().Add(time.Minute))
	cache.advGC()
	if len(cache.store) != 1 {
		t.Error("ttl expired, output store kept untouched")
	}
	mv1, er := cache.Get(k1)
	if er != ErrNotFound {
		t.Error(er)
	}
	if mv1 != nil {
		t.Error("ttl expired, output is not nil")
	}
	mv2, er := cache.Get(k2)
	if mv2 != v2 {
		t.Error("output is not eq input")
	}
}
