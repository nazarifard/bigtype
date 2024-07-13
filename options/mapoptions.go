package options

import "github.com/nazarifard/fastape"

type MapOptions[K comparable, V any] struct {
	hintSize       int
	vTape          fastape.Tape[V]
	kTape          fastape.Tape[K]
	collisionCheck bool
}

func (o *MapOptions[K, V]) WithSize(hintSize int) *MapOptions[K, V] {
	o.hintSize = hintSize
	return o
}

func (o *MapOptions[K, V]) WithVTape(t fastape.Tape[V]) *MapOptions[K, V] {
	o.vTape = t
	return o
}

func (o *MapOptions[K, V]) WithKTape(t fastape.Tape[K]) *MapOptions[K, V] {
	o.kTape = t
	return o
}

func (o *MapOptions[K, V]) Size() int {
	return o.hintSize
}

func (o *MapOptions[K, V]) VTape() fastape.Tape[V] {
	return o.vTape
}
func (o *MapOptions[K, V]) KTape() fastape.Tape[K] {
	return o.kTape
}

func (o *MapOptions[K, V]) WithCollisionCheck(collisionCheck bool) *MapOptions[K, V] {
	o.collisionCheck = collisionCheck
	return o
}

func (o *MapOptions[K, V]) CollisionCheck() bool {
	return o.collisionCheck
}
