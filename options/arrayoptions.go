package options

import (
	"github.com/nazarifard/fastape"
)

type ArrayOptions[V any] struct {
	size         int
	tape         fastape.Tape[V]
	unexpandable bool
}

func (o *ArrayOptions[V]) WithSize(size int) *ArrayOptions[V] {
	o.size = size
	return o
}

func (o *ArrayOptions[V]) WithMarshal(t fastape.Tape[V]) *ArrayOptions[V] {
	o.tape = t
	return o
}

func (o *ArrayOptions[V]) WithExpandable(expandable bool) *ArrayOptions[V] {
	o.unexpandable = !expandable
	return o
}

func (o *ArrayOptions[V]) Size() int {
	return o.size
}

func (o *ArrayOptions[V]) Marshal() fastape.Tape[V] {
	return o.tape
}

func (o *ArrayOptions[V]) Expandable() bool {
	return !o.unexpandable
}
