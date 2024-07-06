package mapiter

type MapIterItem[K comparable, V any] struct {
	next    bool
	mapItem struct {
		Key   K
		Value V
	}
}

type MapIter[K comparable, V any] struct {
	Ch      chan MapIterItem[K, V]
	MapItem struct {
		Key   K
		Value V
	}
}

func (it *MapIter[K, V]) Next() bool {
	n := <-it.Ch
	it.MapItem = n.mapItem
	return n.next
}

func (it *MapIter[K, V]) Value() V {
	return it.MapItem.Value
}

func (it *MapIter[K, V]) Key() K {
	return it.MapItem.Key
}
