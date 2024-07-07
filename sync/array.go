package sync

import (
	"sync"

	"github.com/nazarifard/bigtype/internal/basic"
	"github.com/nazarifard/bigtype/internal/options"
)

const nSubArrays = 512

type array[V any] struct {
	subArrays [nSubArrays]*syncArray[V]
}

func (a *array[V]) Set(i int, v V) {
	a.subArrays[i%nSubArrays].Set(i/nSubArrays, v)
}
func (a *array[V]) Update(i int, updateFn func(old V) (new V)) {
	a.subArrays[i%nSubArrays].Update(i, updateFn)
}
func (a *array[V]) Get(i int) V {
	return a.subArrays[i%nSubArrays].Get(i / nSubArrays)
}

func (a *array[V]) Len() int {
	n := 0
	for i := range a.subArrays {
		n += a.subArrays[i].Len()
	}
	return n
}

func NewArray[V any](ops ...any) Array[V] {
	option := options.ParseArrayOptions[V](ops...)
	a := new(array[V]) //must use new because has non-copiable objects
	for i := range a.subArrays {
		subSize := (option.Size + nSubArrays - 1) / nSubArrays
		a.subArrays[i] = newSyncArray[V](subSize, option.VMarshal, option.Expandable)
	}
	return a
}

func newSyncArray[V any](ops ...any) *syncArray[V] {
	return &syncArray[V]{
		arr:    basic.NewArray[V](ops...),
		mutext: &sync.RWMutex{},
	}
}

type syncArray[V any] struct {
	arr    basic.Array[V]
	mutext *sync.RWMutex
}

func (s *syncArray[V]) Len() int {
	s.mutext.RLock()
	defer s.mutext.RUnlock()
	return s.arr.Len()
}

func (s *syncArray[V]) Set(index int, v V) {
	s.mutext.Lock()
	defer s.mutext.Unlock()
	s.arr.Set(index, v)
}

func (s *syncArray[V]) Update(index int, updateFn func(old V) (new V)) {
	s.mutext.Lock()
	defer s.mutext.Unlock()
	s.arr.Update(index, updateFn)
}

func (s *syncArray[V]) Get(index int) V {
	s.mutext.RLock()
	defer s.mutext.RUnlock()
	return s.arr.Get(index)
}
