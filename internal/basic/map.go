package basic

import (
	"github.com/nazarifard/bigtype/internal/hash"
	"github.com/nazarifard/bigtype/options"
)

type bigMap[K comparable, V any] struct {
	htree  *bigTree[uint64, int] //hash[K]-> index
	keys   Array[K]
	values Array[V]

	collitionMap   map[K]V
	kHash          hash.Hashable[K]
	checkCollition bool
}

func NewMap[K comparable, V any](ops ...any) Map[K, V] {
	opt := ParsMapOptions[K, V](ops...)
	var vo options.ArrayOptions[V]
	var ko options.ArrayOptions[K]
	var ho options.ArrayOptions[int]
	vo.WithSize(opt.Size()).WithExpandable(true).WithMarshal(opt.Marshal())
	ko.WithSize(opt.Size()).WithExpandable(true)
	ho.WithSize(opt.Size()).WithExpandable(true)
	if isNumber[K]() {
		return makeTree[K, V](vo)
	}
	return &bigMap[K, V]{
		htree:          newTree[uint64, int](ho), //TODO hint size
		collitionMap:   make(map[K]V),
		kHash:          hash.NewHash[K](nil),
		keys:           NewArray[K](ko),
		values:         NewArray[V](vo),
		checkCollition: opt.CollisionCheck(),
	}
}

func (m *bigMap[K, V]) Len() int {
	return int(len(m.collitionMap) + m.htree.Len()) // array[0] is not used
}

func (m *bigMap[K, V]) HUpdate(hash uint64, key K, updateFn func(old V) (new V)) {
	index, ok := m.htree.Get(hash)
	if ok { //update
		if !m.checkCollition {
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
		if !m.checkCollition {
			m.keys.Set(index+1, key)
			m.values.Set(index+1, value)
		} else {
			oldKey := m.keys.Get(index + 1)
			if oldKey == key { //just update
				m.keys.Set(index+1, key)
				m.values.Set(index+1, value)
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
	h := m.kHash.Hash(key)
	m.HSet(h, key, value)
}

func (m *bigMap[K, V]) HGet(hash uint64, key K) (value V, ok bool) {
	index, ok := m.htree.Get(hash)
	if !ok {
		return
	}
	if !m.checkCollition {
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
	if m.checkCollition {
		for k, v := range m.collitionMap {
			if next {
				next = f(k, v)
			}
		}
	}
}

func (m *bigMap[K, V]) Delete(key K) {
	hash := m.kHash.Hash(key)
	index, ok := m.htree.Get(hash)
	if ok {
		m.keys.Delete(index)
		m.values.Delete(index)
	}
}
