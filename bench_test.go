package memoryM

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func BenchmarkSet(b *testing.B) {
	manager := New()
	ttl := time.Now().Add(time.Minute)
	keys := make([]string, 1000000)
	for i := 0; i < 1000000; i++ {
		keys[i] = strconv.Itoa(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Set(keys[i%1000000], struct{}{}, ttl)
	}
	b.StopTimer()
}

func BenchmarkGet(b *testing.B) {
	manager := New()
	ttl := time.Now().Add(time.Minute)
	keys := make([]string, 1000000)
	for i := 0; i < 1000000; i++ {
		manager.Set(keys[i], struct{}{}, ttl)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Get(keys[i%1000000])
	}
	b.StopTimer()
}

func BenchmarkSetGet(b *testing.B) {
	manager := New()
	ttl := time.Now().Add(time.Minute)
	keys := make([]string, 1000000)
	for i := 0; i < 1000000; i++ {
		keys[i] = strconv.Itoa(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Set(keys[i%1000000], struct{}{}, ttl)
		manager.Get(keys[i%1000000])
	}
	b.StopTimer()
}

func BenchmarkSimpleTTL(b *testing.B) {
	manager := New()
	n := 5000000
	for i := 0; i < n; i++ {
		off := time.Minute
		if i%1000000 == 0 {
			off = -time.Minute
		}
		ttl := time.Now().Add(off)
		manager.Set(strconv.Itoa(i), struct{}{}, ttl)
	}

	start := time.Now()

	manager.simpleGC()

	fmt.Println(n, "Simple: ", time.Now().Sub(start))
	//r, _ := manager.Get("0")
	//fmt.Println(len(manager.store), r)
	b.SkipNow()
}

func BenchmarkAdvTTL(b *testing.B) {
	manager := New()
	n := 5000000
	for i := 0; i < n; i++ {
		off := time.Minute
		if i%1000000 == 0 {
			off = -time.Minute
		}
		ttl := time.Now().Add(off)
		manager.Set(strconv.Itoa(i), struct{}{}, ttl)
	}

	start := time.Now()

	manager.advGC()

	fmt.Println(n, "Advanced: ", time.Now().Sub(start))
	//r, _ := manager.Get("0")
	//fmt.Println(len(manager.store), r)
	b.SkipNow()
}
