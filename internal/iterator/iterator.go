package iterator

type IterItem[K comparable, V any] struct {
	Key   K
	Value V
	Ok    bool
}

type Iterator[K comparable, V any] struct {
	Ch chan IterItem[K, V]
}

func (it *Iterator[K, V]) Next() (key K, value V, ok bool) {
	item := <-it.Ch
	return item.Key, item.Value, item.Ok
}

// func (it *Iterator[K, V]) Value() V {
// 	return it.iterItem.Value
// }

// func (it *Iterator[K, V]) Key() K {
// 	return it.iterItem.Key
// }
