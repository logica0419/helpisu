package helpisu

import "sync"

/*
SafeCache ジェネリックで、スレッドセーフなマップキャッシュ

Cacheと違い、ロックを取った状態で複数の操作ができます
リセットしても初期キャパシティを記憶しています
*/
type SafeCache[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
	c  int
}

// NewSafeCache 新たなSafeCacheを作成
func NewSafeCache[K comparable, V any](capacity int) *SafeCache[K, V] {
	return &SafeCache[K, V]{
		mu: sync.RWMutex{},
		m:  make(map[K]V, capacity),
		c:  capacity,
	}
}

// Get 指定したKeyのキャッシュを取得
func (c *SafeCache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok = c.m[key]

	return
}

// Set 指定したKey-Valueのセットをキャッシュに入れる
func (c *SafeCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = value
}

// Delete 指定したKeyのキャッシュを削除
func (c *SafeCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.m, key)
}

// Reset 全てのキャッシュを削除
func (c *SafeCache[K, V]) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m = make(map[K]V, c.c)
}

// SafeCacheWithLock ロックが取れた状態のSafeCache
type SafeCacheWithLock[K comparable, V any] struct {
	m map[K]V
	c int
}

// WithLock ロックを取った状態で指定された処理を実行
func (c *SafeCache[K, V]) WithLock(callback func(c *SafeCacheWithLock[K, V]) error) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cl := &SafeCacheWithLock[K, V]{m: c.m, c: c.c}
	if err = callback(cl); err != nil {
		c.m = cl.m
	}

	return
}

// Get 指定したKeyのキャッシュを取得
func (c *SafeCacheWithLock[K, V]) Get(key K) (value V, ok bool) {
	value, ok = c.m[key]

	return
}

// Set 指定したKey-Valueのセットをキャッシュに入れる
func (c *SafeCacheWithLock[K, V]) Set(key K, value V) {
	c.m[key] = value
}

// Delete 指定したKeyのキャッシュを削除
func (c *SafeCacheWithLock[K, V]) Delete(key K) {
	delete(c.m, key)
}

// Reset 全てのキャッシュを削除
func (c *SafeCacheWithLock[K, V]) Reset() {
	c.m = make(map[K]V, c.c)
}
