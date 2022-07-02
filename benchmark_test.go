package helpisu_test

import (
	"testing"

	"github.com/logica0419/helpisu"
)

var Result []string

var cache = helpisu.NewCache[int, int](100)

func BenchmarkCacheNTimes(b *testing.B) {
	b.ReportAllocs()

	cache.Reset()

	for n := 1; n < b.N; n++ {
		cache.Set(n%100, n)
		_, ok := cache.Get(n % 100)

		if !ok {
			b.Errorf("cache.Get(%d) failed", n)
		}

		cache.Delete(n % 100)
	}
}

func BenchmarkCacheNTimesParallel(b *testing.B) {
	b.ReportAllocs()

	cache.Reset()

	for n := 1; n < b.N; n++ {
		go func(n int) {
			cache.Set(n%100, n)
			cache.Get(n % 100)
			cache.Delete(n % 100)
		}(n)
	}
}
