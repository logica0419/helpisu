package helpisu

import "sync"

/*
Cache ジェネリックで、スレッドセーフなマップキャッシュ
	sync.Mapのジェネリックなラッパーです
*/
type Cache[K comparable, V any] struct {
	m *sync.Map
}

// NewCache 新たなCacheを作成
func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		m: &sync.Map{},
	}
}

// Get 指定したKeyのキャッシュを取得
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	cache, ok := c.m.Load(key)
	if !ok {
		return
	}

	value, ok = cache.(V)

	return
}

// GetAndDelete 指定したKeyのキャッシュを取得して削除
func (c *Cache[K, V]) GetAndDelete(key K) (value V, ok bool) {
	cache, ok := c.m.LoadAndDelete(key)
	if !ok {
		return
	}

	value, ok = cache.(V)

	return
}

// Set 指定したKey-Valueのセットをキャッシュに入れる
func (c *Cache[K, V]) Set(key K, value V) {
	c.m.Store(key, value)
}

// Delete 指定したKeyのキャッシュを削除
func (c *Cache[K, V]) Delete(key K) {
	c.m.Delete(key)
}

// Reset 全てのキャッシュを削除
func (c *Cache[K, V]) Reset() {
	c.m = &sync.Map{}
}
