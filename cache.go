package helpisu

import "sync"

type Cache[K comparable, V any] struct {
	sync.RWMutex
	m map[K]V
	c int
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		m: make(map[K]V, capacity),
		c: capacity,
	}
}

func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.RLock()
	value, ok = c.m[key]
	c.RUnlock()
	return
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.Lock()
	c.m[key] = value
	c.Unlock()
}

func (c *Cache[K, V]) Delete(key K) {
	c.Lock()
	delete(c.m, key)
	c.Unlock()
}

func (c *Cache[K, V]) Reset() {
	c.Lock()
	c.m = make(map[K]V, c.c)
	c.Unlock()
}
