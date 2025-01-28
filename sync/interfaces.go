// Developer: bahador.nazarifard@gmail.com

package sync

type Array[V any] interface {
	Set(index int, v V)
	Get(index int) V
	Len() int
	Update(index int, updateFn func(old V) (new V))
	Seq(yield func(V) bool)
}

type Updatable[K comparable, V any] interface {
	Update(key K, updateFn func(old V) (new V))
}

type Map[K comparable, V any] interface {
	Updatable[K, V]
	Len() int
	Set(key K, value V)
	Get(key K) (value V, ok bool)
	SetMany(items map[K]V)
	Seq(f func(Key K, Value V) bool)
	Delete(key K)
}

//type tree[K kNumber, V any] Map[K, V]

type kNumber interface {
	uint | ~int | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// type hashable[K comparable] interface {
// 	Hash(K) uint64
// }

// type hMap[H hashable[K], K comparable, V any] interface {
// 	Hasher() H
// 	Map[K, V]
// }
