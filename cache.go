package helpisu

import "sync"

/*
Cache ジェネリックで、スレッドセーフなマップキャッシュ

リセットしても初期キャパシティを記憶しています
*/
type Cache[K comparable, V any] struct {
	sync.RWMutex
	m map[K]V
	c int
}

// NewCache 新たなCacheを作成
func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		m: make(map[K]V, capacity),
		c: capacity,
	}
}

// Get 指定したKeyのキャッシュを取得
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.RLock()
	value, ok = c.m[key]
	c.RUnlock()

	return
}

// Set 指定したKey-Valueのセットをキャッシュに入れる
func (c *Cache[K, V]) Set(key K, value V) {
	c.Lock()
	c.m[key] = value
	c.Unlock()
}

// Delete 指定したKeyのキャッシュを削除
func (c *Cache[K, V]) Delete(key K) {
	c.Lock()
	delete(c.m, key)
	c.Unlock()
}

// Reset 全てのキャッシュを削除
func (c *Cache[K, V]) Reset() {
	c.Lock()
	c.m = make(map[K]V, c.c)
	c.Unlock()
}
