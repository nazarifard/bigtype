package basic

import (
	"sync"
)

type SyncMap[K comparable, V any] struct {
	m     Map[K, V]
	mutex *sync.RWMutex
}

func NewSyncMap[K comparable, V any](ops ...any) *SyncMap[K, V] {
	m := SyncMap[K, V]{
		m:     NewMap[K, V](ops...), //.(*bigMap1[K, V]),
		mutex: &sync.RWMutex{},
	}
	//m.m.iterator = nil
	return &m
}

func (m *SyncMap[K, V]) Len() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.m.Len()
}

func (m *SyncMap[K, V]) HSet(hash uint64, key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.m.HSet(hash, key, value)
}

func (m *SyncMap[K, V]) HGet(hash uint64, key K) (value V, ok bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.m.HGet(hash, key)
}

func (m *SyncMap[K, V]) HUpdate(hash uint64, key K, updateFn func(old V) (new V)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.m.HUpdate(hash, key, updateFn)
}

func (m *SyncMap[K, V]) Update(key K, updateFn func(old V) (new V)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.m.Update(key, updateFn)
}

func (m *SyncMap[K, V]) Set(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.m.Set(key, value)
}

func (m *SyncMap[K, V]) Get(key K) (value V, ok bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.m.Get(key)
}

func (m *SyncMap[K, V]) Seq(f func(Key K, Value V) bool) {
	m.mutex.RLock()
	m.m.Seq(f)
	m.mutex.RUnlock()
}

func (m *SyncMap[K, V]) SetMany(items map[K]V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.m.SetMany(items)
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.mutex.RLock()
	m.m.Delete(key)
	m.mutex.RUnlock()
}
