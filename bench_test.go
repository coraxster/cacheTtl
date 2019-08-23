package cacheTtl

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func BenchmarkSet(b *testing.B) {
	cache := New()
	ttl := time.Now().Add(time.Minute)
	keys := make([]string, 1000000)
	for i := 0; i < 1000000; i++ {
		keys[i] = strconv.Itoa(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(keys[i%1000000], struct{}{}, ttl)
	}
	b.StopTimer()
}

func BenchmarkGet(b *testing.B) {
	cache := New()
	ttl := time.Now().Add(time.Minute)
	keys := make([]string, 1000000)
	for i := 0; i < 1000000; i++ {
		cache.Set(keys[i], struct{}{}, ttl)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(keys[i%1000000])
	}
	b.StopTimer()
}

func BenchmarkSetGet(b *testing.B) {
	cache := New()
	ttl := time.Now().Add(time.Minute)
	keys := make([]string, 1000000)
	for i := 0; i < 1000000; i++ {
		keys[i] = strconv.Itoa(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n := i % 1000000
		cache.Set(keys[n], struct{}{}, ttl)
		cache.Get(keys[n])
	}
	b.StopTimer()
}

func BenchmarkSimpleTTL(b *testing.B) {
	cache := New()
	n := 5000000
	for i := 0; i < n; i++ {
		off := time.Minute
		if i%100000 == 0 {
			off = -time.Minute
		}
		ttl := time.Now().Add(off)
		cache.Set(strconv.Itoa(i), struct{}{}, ttl)
	}

	start := time.Now()
	b.ResetTimer()
	cache.simpleGC()
	fmt.Println(n, "SimpleGC: ", time.Now().Sub(start))
	b.SkipNow()
}

func BenchmarkAdvTTL(b *testing.B) {
	cache := New()
	n := 5000000
	expired := 0
	for i := 0; i < n; i++ {
		off := time.Minute
		if i%100000 == 0 {
			off = -time.Minute
			expired++
		}
		ttl := time.Now().Add(off)
		cache.Set(strconv.Itoa(i), struct{}{}, ttl)
	}
	start := time.Now()
	b.ResetTimer()
	cache.advGC()
	fmt.Println(n, "AdvancedGC: ", time.Now().Sub(start))
	fmt.Println(len(cache.store), n-expired)
	if len(cache.store) != n-expired {
		b.Error("GC worked not so good")
	}
	b.SkipNow()
}
