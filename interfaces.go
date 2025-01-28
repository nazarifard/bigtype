// Developer: bahador.nazarifard@gmail.com

package bigtype

import "iter"

type Array[V any] interface {
	Set(index int, v V)
	Get(index int) V
	Len() int
	Update(index int, updateFn func(old V) (new V))
	Delete(index int)
	All() iter.Seq[V]
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
	Range(f func(Key K, Value V) bool)
	Delete(key K)
	All() iter.Seq2[K, V]
}
