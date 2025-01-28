// Developer: bahador.nazarifard@gmail.com

package bigtype

type Array[V any] interface {
	Set(index int, v V)
	Get(index int) V
	Len() int
	Update(index int, updateFn func(old V) (new V))
	Delete(index int)
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
	Seq(yield func(K, V) bool)
	Delete(key K)
}
