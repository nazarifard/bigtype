package basic

import (
	"fmt"
	"os"
	"reflect"

	"github.com/nazarifard/bigtype/internal/hash"
	"github.com/nazarifard/bigtype/options"
)

type node[H kNumber] struct {
	hkey  H
	left  uint32
	right uint32
}

type bigTree[K kNumber, V any] struct {
	keys       Array[K]
	values     Array[V]
	indexTable Array[node[K]] //*bigFixedArray[node[K]]
	root       uint32
	free       uint32
	hasher     hash.Hashable[K]
}

func newTree[K kNumber, V any](ops ...any) *bigTree[K, V] {
	if reflect.ValueOf(*new(K)).Kind() == reflect.String {
		fmt.Println("Error: bigtype.bigTree does not support string keys. use bigtype.Map")
		os.Exit(1)
	}
	opt := ParsTreeOptions[K, V](ops...)

	var vopt options.ArrayOptions[V]
	vopt.WithSize(opt.Size() + 1).WithMarshal(opt.Marshal())

	return &bigTree[K, V]{
		indexTable: NewFixedArray[node[K]](opt.Size()+1, true),
		keys:       NewFixedArray[K](opt.Size()+1, true),
		values:     NewArray[V](vopt),
		hasher:     hash.NewHash[K](),
		root:       0,
		free:       1,
	}
}

func (t *bigTree[K, V]) Len() int {
	return int(t.free - 1) // array[0] is not used
}

func (t *bigTree[K, V]) SetMany(m map[K]V) {
	for k, v := range m {
		t.Set(K(t.hasher.Hash(k)), v)
	}
}

func (t *bigTree[K, V]) Set(key K, value V) {
	index := t.insertIterative(t.root, K(t.hasher.Hash(key)))
	t.values.Set(int(index), value)
	t.keys.Set(int(index), key)
}

func (t *bigTree[K, V]) Update(key K, updateFn func(old V) (new V)) {
	index := t.insertIterative(t.root, K(t.hasher.Hash(key)))
	old := t.values.Get(int(index))
	t.values.Set(int(index), updateFn(old))
}

func (t *bigTree[K, V]) insertIterative(root uint32, hkey K) (index uint32) {
	const NULL uint32 = 0

	newNode := node[K]{hkey: hkey}
	if root == NULL {
		t.indexTable.Set(int(t.free), newNode)
		t.root = t.free
		t.free++
		return t.root
	}

	parent := NULL
	curr := root

	for curr != NULL {
		parent = curr
		currentNode := t.indexTable.Get(int(curr))
		if hkey == currentNode.hkey {
			return curr
		} else if hkey < currentNode.hkey {
			curr = currentNode.left
		} else {
			curr = currentNode.right
		}
	}

	parentNode := t.indexTable.Get(int(parent))
	if hkey < parentNode.hkey {
		parentNode.left = t.free
	} else {
		parentNode.right = t.free
	}
	t.indexTable.Set(int(parent), parentNode)
	t.indexTable.Set(int(t.free), newNode)
	t.free++
	return t.free - 1
}

func (t *bigTree[K, V]) Get(key K) (value V, ok bool) {
	index, ok := t.iterativeSearch(t.root, K(t.hasher.Hash(key)))
	if ok {
		return t.values.Get(int(index)), true
	}
	return
}

func (t *bigTree[K, V]) iterativeSearch(root uint32, hkey K) (index uint32, ok bool) {
	for root != 0 {
		node := t.indexTable.Get(int(root))
		if hkey > node.hkey {
			root = node.right
		} else if hkey < node.hkey {
			root = node.left
		} else {
			return root, true
		}
	}
	return 0, false
}

func (t *bigTree[K, V]) inOrder(root uint32) {
	if root != 0 {
		t.inOrder(t.indexTable.Get(int(root)).left)
		hash := t.indexTable.Get(int(root))
		fmt.Printf("index:%d, hash:%v node:{%v, %v}\n", root, hash.hkey, t.keys.Get(int(root)), t.values.Get(int(root)))
		t.inOrder(t.indexTable.Get(int(root)).right)
	}
}

func (t *bigTree[K, V]) Range(f func(key K, value V) bool) {
	next := true
	for i := 1; i < int(t.free) && next; i++ {
		value := t.values.Get(i)
		key := t.keys.Get(i)
		next = f(key, value)
	}
}

func (t *bigTree[K, V]) Delete(key K) {
	//TODO index should be deleted from tree also
	index := t.insertIterative(t.root, K(t.hasher.Hash(key)))
	t.values.Delete(int(index))
	t.keys.Delete(int(index))
}

// TODO not implemented yet
func (t *bigTree[K, V]) HSet(hash uint64, key K, value V)                         {}
func (t *bigTree[K, V]) HGet(hash uint64, key K) (value V, ok bool)               { return }
func (t *bigTree[K, V]) HUpdate(hash uint64, key K, updateFn func(old V) (new V)) {}
