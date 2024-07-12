package options

import (
	"github.com/nazarifard/fastape"
)

type TreeOptions[K comparable, V any] struct {
	hintSize int
	tape     fastape.Tape[V]
	// expandable bool
}

func (o *TreeOptions[K, V]) WithSize(hintSize int) *TreeOptions[K, V] {
	o.hintSize = hintSize
	return o
}

func (o *TreeOptions[K, V]) WithMarshal(t fastape.Tape[V]) *TreeOptions[K, V] {
	o.tape = t
	return o
}

// func (o *TreeOptions[K,V]) WithExtandable(extandable bool) *TreeOptions[K,V] {
// 	o.expandable = extandable
// 	return o
// }

func (o *TreeOptions[K, V]) Size() int {
	return o.hintSize
}

func (o *TreeOptions[K, V]) Marshal() fastape.Tape[V] {
	return o.tape
}

// func (o *TreeOptions[K,V]) Expandable() bool {
// 	return o.expandable
// }
