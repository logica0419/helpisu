package helpisu

import "sync"

/*
Cache ジェネリックで、スレッドセーフなマップキャッシュ
	リセットしても初期キャパシティを記憶しています
*/
type Cache[K comparable, V any] struct {
	m *sync.Pool
	c int
}

// NewCache 新たなCacheを作成
func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		m: &sync.Pool{
			New: func() interface{} {
				return make(map[K]V, capacity)
			},
		},
		c: capacity,
	}
}

// Get 指定したKeyのキャッシュを取得
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	cache, _ := c.m.Get().(map[K]V)
	defer c.m.Put(cache)
	value, ok = cache[key]

	return
}

// Set 指定したKey-Valueのセットをキャッシュに入れる
func (c *Cache[K, V]) Set(key K, value V) {
	cache, _ := c.m.Get().(map[K]V)
	cache[key] = value
	c.m.Put(cache)
}

// Delete 指定したKeyのキャッシュを削除
func (c *Cache[K, V]) Delete(key K) {
	cache, _ := c.m.Get().(map[K]V)
	delete(cache, key)
	c.m.Put(cache)
}

// Reset 全てのキャッシュを削除
func (c *Cache[K, V]) Reset() {
	c.m.Put(make(map[K]V, c.c))
}
