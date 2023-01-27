package cache

import (
	"sync"
)

type resetter interface {
	Reset()
}

var generatedCaches = []resetter{}

/*
Cache ジェネリックで、スレッドセーフなマップキャッシュ

	sync.Mapのジェネリックなラッパーです
*/
type Cache[K comparable, V any] struct {
	m *sync.Map
	s func(key K, value V)
	d func(key K)
	r func()
}

// New 新たなCacheを作成
func New[K comparable, V any]() *Cache[K, V] {
	c := Cache[K, V]{
		m: &sync.Map{},
		s: nil,
		d: nil,
		r: nil,
	}

	generatedCaches = append(generatedCaches, &c)

	return &c
}

// Get 指定したKeyのキャッシュを取得
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return
	}

	value, ok = v.(V)

	return
}

// GetAndDelete 指定したKeyのキャッシュを取得して削除
func (c *Cache[K, V]) GetAndDelete(key K) (value V, ok bool) {
	v, ok := c.m.LoadAndDelete(key)
	if !ok {
		return
	}

	value, ok = v.(V)

	return
}

// Set 指定したKey-Valueのセットをキャッシュに入れる
func (c *Cache[K, V]) Set(key K, value V) {
	c.m.Store(key, value)

	if c.s != nil {
		c.s(key, value)
	}
}

// Delete 指定したKeyのキャッシュを削除
func (c *Cache[K, V]) Delete(key K) {
	c.m.Delete(key)

	if c.d != nil {
		c.d(key)
	}
}

// ForEach キャッシュの全ての要素に対して処理を行う
func (c *Cache[K, V]) ForEach(f func(key K, value V) error) (err error) {
	c.m.Range(func(key, value interface{}) bool {
		k, _ := key.(K)
		v, _ := value.(V)

		err = f(k, v)
		return err == nil
	})

	return
}

// Reset 全てのキャッシュを削除
func (c *Cache[K, V]) Reset() {
	c.m = &sync.Map{}

	if c.r != nil {
		c.r()
	}
}

// ResetAll `NewCache()`で生成した全てのキャッシュをリセット
func ResetAll() {
	for _, c := range generatedCaches {
		c.Reset()
	}
}
