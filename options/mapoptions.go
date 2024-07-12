package options

import "github.com/nazarifard/fastape"

type MapOptions[K comparable, V any] struct {
	hintSize int
	vTape    fastape.Tape[V]
	// expandable bool
	collisionCheck bool
}

func (o *MapOptions[K, V]) WithSize(hintSize int) *MapOptions[K, V] {
	o.hintSize = hintSize
	return o
}

func (o *MapOptions[K, V]) WithMarshal(t fastape.Tape[V]) *MapOptions[K, V] {
	o.vTape = t
	return o
}

func (o *MapOptions[K, V]) Size() int {
	return o.hintSize
}

func (o *MapOptions[K, V]) Marshal() fastape.Tape[V] {
	return o.vTape
}

// func (o *MapOptions[K,V]) WithExtandable(extandable bool) *MapOptions[K,V] {
// 	o.expandable = extandable
// 	return o
// }

// func (o *MapOptions[K,V]) Expandable() bool {
// 	return o.expandable
// }

func (o *MapOptions[K, V]) WithCollisionCheck(collisionCheck bool) *MapOptions[K, V] {
	o.collisionCheck = collisionCheck
	return o
}

func (o *MapOptions[K, V]) CollisionCheck() bool {
	return o.collisionCheck
}
