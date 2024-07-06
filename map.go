package bigtype

import (
	"sync"

	"github.com/nazarifard/bigtype/internal/basic"
	"github.com/nazarifard/bigtype/internal/hash"
	"github.com/nazarifard/bigtype/internal/options"
)

const nSubMaps = 512

type mapIterItem[K comparable, V any] struct {
	k  K
	v  V
	ok bool
}

type bigMap[K comparable, V any] struct {
	subMaps [nSubMaps]*basic.SyncMap[K, V]
	kHash   hash.Hashable[K]
	Chan    *chan mapIterItem[K, V]
	mutexCh *sync.RWMutex
}

func NewMap[K comparable, V any](ops ...any) Map[K, V] {
	if isNumber[K]() {
		return makeTree[K, V](ops...)
	}
	option := options.ParseMapOptions[K, V](ops...)
	var newOps []any
	if len(ops) > 0 {
		hintSize := (option.HintSize + nSubMaps - 1) / nSubMaps
		newOps = append(newOps, hintSize)
		if len(ops) > 1 {
			newOps = append(newOps, ops[1:])
		}
	}

	var m bigMap[K, V]
	m.kHash = hash.NewHash[K](option.KMarshal)
	m.mutexCh = &sync.RWMutex{}
	for i := range m.subMaps {
		m.subMaps[i] = basic.NewSyncMap[K, V](newOps...)
	}
	return &m
}

func (m *bigMap[K, V]) SetMany(in map[K]V) {
	for k, v := range in {
		m.Set(k, v)
	}
}

func (m *bigMap[K, V]) Set(key K, value V) {
	hash := m.kHash.Hash(key)
	m.subMaps[hash%nSubMaps].HSet(hash, key, value)
}

func (m *bigMap[K, V]) Update(key K, updateFn func(old V) (new V)) {
	hash := m.kHash.Hash(key)
	m.subMaps[hash%nSubMaps].HUpdate(hash, key, updateFn)
}

func (m *bigMap[K, V]) Get(key K) (value V, ok bool) {
	hash := m.kHash.Hash(key)
	return m.subMaps[hash%nSubMaps].HGet(hash, key)
}

func (m *bigMap[K, V]) Len() int {
	n := 0
	for i := range m.subMaps {
		n += m.subMaps[i].Len()
	}
	return n
}

// Just one writer but multiple reader
// but one writer should collect all subtrees data from multi channel
// then push items to new channel
// read from multiple channel.......but push to single channel
func (m *bigMap[K, V]) Range(f func(Key K, Value V) bool) {
	var ch *chan mapIterItem[K, V]
	m.mutexCh.RLock()
	ch = m.Chan
	m.mutexCh.RUnlock()
	if ch == nil {
		ch = func() *chan mapIterItem[K, V] {
			ch := make(chan mapIterItem[K, V])
			go func() {
				defer func() {
					close(ch)
					m.mutexCh.Lock()
					m.Chan = nil
					m.mutexCh.Unlock()
				}()
				for i := range m.subMaps {
					m.subMaps[i].Range(func(key K, value V) bool {
						ch <- mapIterItem[K, V]{k: key, v: value, ok: true}
						return true
					})
				}
			}()
			return &ch
		}()
		m.mutexCh.Lock()
		m.Chan = ch
		m.mutexCh.Unlock()
	}
	for item := range *ch {
		if !f(item.k, item.v) {
			break
		}
	}
}
