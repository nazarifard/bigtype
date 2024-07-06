package basic

import (
	"github.com/nazarifard/bigtype/internal/hash"
	"github.com/nazarifard/bigtype/internal/options"
)

type bigMap[K comparable, V any] struct {
	htree  *bigTree[uint64, int] //hash[K]-> index
	keys   Array[K]
	values Array[V]

	collitionMap   map[K]V
	kHash          hash.Hashable[K]
	CheckCollition bool
}

func NewMap[K comparable, V any](ops ...any) Map[K, V] {
	if isNumber[K]() {
		return makeTree[K, V]()
	}
	option := options.ParseMapOptions[K, V](ops...)
	return &bigMap[K, V]{
		htree:          newTree[uint64, int](option.HintSize, nil, true), //TODO hint size
		collitionMap:   make(map[K]V),
		kHash:          hash.NewHash[K](option.KMarshal),
		keys:           NewArray[K](option.HintSize, option.KMarshal, true),
		values:         NewArray[V](option.HintSize, option.VMarshal, true),
		CheckCollition: option.CheckCollition,
	}
}

func (m *bigMap[K, V]) Len() int {
	return int(len(m.collitionMap) + m.htree.Len()) // array[0] is not used
}

func (m *bigMap[K, V]) HUpdate(hash uint64, key K, updateFn func(old V) (new V)) {
	index, ok := m.htree.Get(hash)
	if ok { //update
		if !m.CheckCollition {
			//m.keys.Set(index, key)
			m.values.Update(index, updateFn)
		} else {
			oldKey := m.keys.Get(index)
			if oldKey == key { //just update
				//m.keys.Set(index, key)
				m.values.Update(index, updateFn)
			} else {
				//real collicion
				m.collitionMap[key] = updateFn(m.collitionMap[key])
			}
		}
	} else { //new insert
		Len := m.htree.Len()
		m.htree.Set(hash, Len) //index+1
		m.keys.Set(Len+1, key)
		m.values.Set(Len+1, updateFn(*new(V)))
	}
}
func (m *bigMap[K, V]) Update(key K, updateFn func(old V) (new V)) {
	m.HUpdate(m.kHash.Hash(key), key, updateFn)
}

func (m *bigMap[K, V]) HSet(hash uint64, key K, value V) {
	index, ok := m.htree.Get(hash)
	if ok { //update
		if !m.CheckCollition {
			m.keys.Set(index, key)
			m.values.Set(index, value)
		} else {
			oldKey := m.keys.Get(index)
			if oldKey == key { //just update
				m.keys.Set(index, key)
				m.values.Set(index, value)
			} else {
				//real collicion
				m.collitionMap[key] = value
			}
		}
	} else { //new insert
		Len := m.htree.Len()
		m.htree.Set(hash, Len) //index+1
		m.keys.Set(Len+1, key)
		m.values.Set(Len+1, value)
	}
}

func (m *bigMap[K, V]) Set(key K, value V) {
	m.HSet(m.kHash.Hash(key), key, value)
}

func (m *bigMap[K, V]) HGet(hash uint64, key K) (value V, ok bool) {
	index, ok := m.htree.Get(hash)
	if !ok {
		return
	}
	if !m.CheckCollition {
		//m.keys.Set(index, key)
		value = m.values.Get(index)
		return
	}

	//first check collitionMap
	value, ok = m.collitionMap[key]
	if !ok {
		value = m.values.Get(index)
	}
	return
}

func (m *bigMap[K, V]) Get(key K) (value V, ok bool) {
	return m.HGet(m.kHash.Hash(key), key)
}

func (m *bigMap[K, V]) SetMany(items map[K]V) {
	for k, v := range items {
		m.Set(k, v)
	}
}

func (m *bigMap[K, V]) Range(f func(key K, value V) bool) {
	next := true
	for i := 1; next && i <= m.htree.Len(); i++ {
		next = f(m.keys.Get(i), m.values.Get(i))
	}
	if m.CheckCollition {
		for k, v := range m.collitionMap {
			if next {
				next = f(k, v)
			}
		}
	}
}
