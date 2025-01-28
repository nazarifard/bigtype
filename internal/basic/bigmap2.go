package basic

import (
	"github.com/nazarifard/bigtype/internal/hash"
	"github.com/nazarifard/bigtype/internal/kv"
	"github.com/nazarifard/fastape"
)

// none of K or V is fixedtype, both have variable size
type bigMap2[K comparable, V any] struct {
	htree *bigTree[uint64, int] //hash[K]-> index
	kv    Array[kv.KV[K, V]]

	collitionMap   map[K]V
	kHash          hash.Hashable[K]
	checkCollition bool
	tape           kv.TapeKV[K, V] //fastape.Tape[kv.KV[K, V]]
	vTape          fastape.Tape[V]
}

func (m *bigMap2[K, V]) Len() int {
	return int(len(m.collitionMap) + m.htree.Len()) // array[0] is not used
}

func (m *bigMap2[K, V]) HUpdate(hash uint64, key K, fn func(V) V) {
	newFn := func(old kv.KV[K, V]) (new kv.KV[K, V]) {
		new.Value = fn(old.Value)
		return
	}
	index, ok := m.htree.Get(hash)
	if ok { //update
		if !m.checkCollition {
			//m.keys.Set(index, key)
			m.kv.Update(index, newFn)
		} else {
			old := m.kv.Get(index)
			if old.Key == key { //just update
				//m.keys.Set(index, key)
				m.kv.Update(index, newFn)
			} else {
				//real collicion
				m.collitionMap[key] = fn(m.collitionMap[key])
			}
		}
	} //else { //key not found
	// Len := m.htree.Len()
	// m.htree.Set(hash, Len)
	// m.kv.Set(Len, kv.KV[K, V]{Key: key, Value: (*new(V))})
	//}
}
func (m *bigMap2[K, V]) Update(key K, updateFn func(old V) (new V)) {
	m.HUpdate(m.kHash.Hash(key), key, updateFn)
}

func (m *bigMap2[K, V]) HSet(hash uint64, key K, value V) {
	item := kv.KV[K, V]{Key: key, Value: value}
	index, ok := m.htree.Get(hash)
	if ok { //update
		if !m.checkCollition {
			m.kv.Set(index, item)
		} else {
			old := m.kv.Get(index)
			if old.Key == key { //just update
				m.kv.Set(index, item)
			} else {
				//real collicion
				m.collitionMap[key] = value
			}
		}
	} else { //new insert
		Len := m.htree.Len()
		m.htree.Set(hash, Len) //index+1
		m.kv.Set(Len, item)
	}
}

func (m *bigMap2[K, V]) Set(key K, value V) {
	h := m.kHash.Hash(key)
	m.HSet(h, key, value)
}

func (m *bigMap2[K, V]) HGet(hash uint64, key K) (value V, ok bool) {
	index, ok := m.htree.Get(hash)
	if !ok {
		return
	}
	origArray, ok := (m.kv).(*array[kv.KV[K, V]])
	if !ok {
		return value, false
	}

	if !m.checkCollition {
		rawItem := origArray.BytesArray.UnsafePtr(index)
		m.tape.Value(rawItem)
	}

	//first check collitionMap
	value, ok = m.collitionMap[key]
	if !ok {
		rawItem := origArray.BytesArray.UnsafePtr(index)
		return m.tape.Value(rawItem)
	}
	return
}

func (m *bigMap2[K, V]) Get(key K) (value V, ok bool) {
	return m.HGet(m.kHash.Hash(key), key)
}

func (m *bigMap2[K, V]) SetMany(items map[K]V) {
	for k, v := range items {
		m.Set(k, v)
	}
}

func (m *bigMap2[K, V]) Seq(f func(K, V) bool) {
	next := true
	for i := 0; next && i < m.htree.Len(); i++ {
		item := m.kv.Get(i)
		next = f(item.Key, item.Value)
	}
	if m.checkCollition {
		for k, v := range m.collitionMap {
			if next {
				next = f(k, v)
			}
		}
	}
}

func (m *bigMap2[K, V]) Delete(key K) {
	hash := m.kHash.Hash(key)
	index, ok := m.htree.Get(hash)
	if ok {
		m.kv.Delete(index)
	}
}

// func (m *bigMap2[K, V]) UnsafeGetValue(index int) (value V, ok bool) {
// 	a, ok := (m.kv).(*array[kv.KV[K, V]])
// 	if ok {
// 		rawItem := a.BytesArray.UnsafePtr(index)
// 		if rawItem == nil {
// 			return value, false
// 		}
// 		//m.vTape.Unroll(rawItem, &value)
// 		lenTape := fastape.LenTape{}
// 		var size int
// 		n, _ := lenTape.Unroll(rawItem, &size)
// 		bs := rawItem[n : n+size]
// 		value = any(bs).(V)
// 		return value, true
// 		// if value, ok = any(bs).(V); !ok { //[]byte
// 		// 	if value, ok = any(unsafe.String(&bs[0], size)).(V); !ok { //string
// 		// 		_, err := m.vTape.Unroll(bs, &value) //unknown
// 		// 		if err != nil {
// 		// 			return *new(V), false
// 		// 		}
// 		// 	}
// 		// }
// 		// return value, true
// 	}
// 	return value, true

// }
