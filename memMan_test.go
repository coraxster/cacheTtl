package memMan

import (
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	k := "key"
	v := "value"
	manager := New()
	manager.Set(k, v, time.Now().Add(time.Minute))
	mv, er := manager.Get(k)
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
	manager := New()
	manager.Set(k, v, time.Now().Add(time.Minute))
	manager.Del(k)
	mv, er := manager.Get(k)
	if er.Error() != "not found" {
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
	manager := New()
	manager.Set(k1, v1, time.Now().Add(-time.Minute))
	manager.Set(k2, v2, time.Now().Add(time.Minute))
	mv1, er := manager.Get(k1)
	if er.Error() != "not found" {
		t.Error(er)
	}
	if mv1 != nil {
		t.Error("ttl expired, output is not nil")
	}
	mv2, er := manager.Get(k2)
	if mv2 != v2 {
		t.Error("output is not eq input")
	}
}

func TestTtlByTime(t *testing.T) {
	k1 := "key1"
	v1 := "value"
	k2 := "key2"
	v2 := "value"
	manager := New()
	manager.Set(k1, v1, time.Now().Add(-time.Minute))
	manager.Set(k2, v2, time.Now().Add(time.Minute))
	ti := make(chan time.Time)
	ticker = ti
	ti <- time.Now()
	if len(manager.store) != 1 {
		t.Error("ttl expired, output store kept untuched")
	}
	mv1, er := manager.Get(k1)
	if er.Error() != "not found" {
		t.Error(er)
	}
	if mv1 != nil {
		t.Error("ttl expired, output is not nil")
	}
	mv2, er := manager.Get(k2)
	if mv2 != v2 {
		t.Error("output is not eq input")
	}
}
