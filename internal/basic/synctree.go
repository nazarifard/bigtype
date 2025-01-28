package basic

import (
	"sync"
)

type SyncTree[K kNumber, V any] struct {
	t     *bigTree[K, V]
	mutex *sync.RWMutex
}

func NewSyncTree[K kNumber, V any](ops ...any) *SyncTree[K, V] {
	t := SyncTree[K, V]{
		t:     newTree[K, V](ops...),
		mutex: &sync.RWMutex{},
	}
	return &t
}

func (t *SyncTree[K, V]) Len() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.t.Len()
}

func (t *SyncTree[K, V]) Set(key K, value V) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.t.Set(key, value)
}

func (t *SyncTree[K, V]) Update(key K, updateFn func(old V) (new V)) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.t.Update(key, updateFn)
}

func (t *SyncTree[K, V]) Get(key K) (value V, ok bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.t.Get(key)
}

func (t *SyncTree[K, V]) Seq(f func(key K, value V) bool) {
	t.mutex.RLock()
	t.t.Seq(f)
	t.mutex.RUnlock()
}

func (t *SyncTree[K, V]) SetMany(items map[K]V) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.t.SetMany(items)
}

func (t *SyncTree[K, V]) Delete(key K) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	t.t.Delete(key)
}
