package sync

import (
	"sync"

	"github.com/nazarifard/bigtype/internal/basic"
)

const nSubTrees = 512

type treeIterItem[K kNumber, V any] struct {
	k  K
	v  V
	ok bool
}

type bigTree[K kNumber, V any] struct {
	subTrees [nSubTrees]*basic.SyncTree[K, V]
	Chan     *chan treeIterItem[K, V]
	mutexCh  *sync.RWMutex
}

func newTree[K kNumber, V any](ops ...any) *bigTree[K, V] {
	options := basic.ParsTreeOptions[K, V](ops...)
	var t bigTree[K, V]
	hintSize := (options.Size() + nSubTrees - 1) / nSubTrees
	options.WithSize(hintSize)
	for i := range t.subTrees {
		t.subTrees[i] = basic.NewSyncTree[K, V](options)
	}
	t.mutexCh = &sync.RWMutex{}
	return &t
}

func (t *bigTree[K, V]) SetMany(m map[K]V) {
	for k, v := range m {
		t.Set(k, v)
	}
}

func (t *bigTree[K, V]) Set(key K, value V) {
	t.subTrees[int(key)%nSubTrees].Set(key, value)
}

func (t *bigTree[K, V]) Update(key K, updateFn func(old V) (new V)) {
	t.subTrees[int(key)%nSubTrees].Update(key, updateFn)
}

func (t *bigTree[K, V]) Get(key K) (value V, ok bool) {
	return t.subTrees[int(key)%nSubTrees].Get(key)
}

func (t *bigTree[K, V]) Len() int {
	n := 0
	for i := range t.subTrees {
		n += t.subTrees[i].Len()
	}
	return n
}

// Just one writer but multiple reader
// but one writer should collect all subtrees data from multi channel
// then push items to new channel
// read from multiple channel.......but push to single channel
func (t *bigTree[K, V]) Seq(f func(Key K, Value V) bool) {
	var ch *chan treeIterItem[K, V]
	t.mutexCh.RLock()
	ch = t.Chan
	t.mutexCh.RUnlock()
	if ch == nil {
		ch = func() *chan treeIterItem[K, V] {
			ch := make(chan treeIterItem[K, V])
			go func() {
				defer func() {
					close(ch)
					t.mutexCh.Lock()
					t.Chan = nil
					t.mutexCh.Unlock()
				}()
				for i := range t.subTrees {
					t.subTrees[i].Seq(func(key K, value V) bool {
						ch <- treeIterItem[K, V]{k: key, v: value, ok: true}
						return true
					})
				}
			}()
			return &ch
		}()
		t.mutexCh.Lock()
		t.Chan = ch
		t.mutexCh.Unlock()
	}
	for item := range *ch {
		if !f(item.k, item.v) {
			break
		}
	}
}

func (t *bigTree[K, V]) Delete(key K) {
	t.subTrees[int(key)%nSubTrees].Delete(key)
}
