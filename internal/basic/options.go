package basic

import (
	"errors"

	marshal "github.com/nazarifard/marshaltap"
)

type ArrayOptions[V any] struct {
	Size       int
	marshal    marshal.Interface[V]
	cloneFn    func(V) V
	expandable bool
}

func (ao *ArrayOptions[V]) WithSize(size int) *ArrayOptions[V] {
	ao.Size = size
	return ao
}

func (ao *ArrayOptions[V]) WithMarshal(m marshal.Interface[V]) *ArrayOptions[V] {
	ao.marshal = m
	return ao
}

func (ao *ArrayOptions[V]) WithCloneFn(fn func(V) V) *ArrayOptions[V] {
	ao.cloneFn = fn
	return ao
}

func (ao *ArrayOptions[V]) WithExtandable(extandable bool) *ArrayOptions[V] {
	ao.expandable = extandable
	return ao
}

type TreeOptions[V any] struct {
	ArrayOptions[V]
}

type MapOptions[K comparable, V any] struct {
	ArrayOptions[V]
	checkCollition bool
}

func (mo *MapOptions[K, V]) WithCheckColission(checkCollision bool) *MapOptions[K, V] {
	mo.checkCollition = checkCollision
	return mo
}

func ParsArrayOptions[V any](ops ...any) (ao ArrayOptions[V]) {
	ao.expandable = true //by default
	if len(ops) == 1 {
		switch t := ops[0].(type) {
		case int:
			ao.WithSize(t).WithExtandable(false)
		case ArrayOptions[V]:
			ao = t
		default:
			panic(errors.New("invalid array options"))
		}
	}
	return ao
}

func ParsMapOptions[K comparable, V any](ops ...any) (mo MapOptions[K, V]) {
	var v V
	mo.expandable = true //by default
	if !IsRequiredToMarshal(v) {
		switch len(ops) {
		case 0:
			return mo
		case 1:
			size, ok := ops[0].(int)
			if ok {
				mo.WithSize(size).WithExtandable(true)
				return mo
			}

			mo, ok := ops[0].(MapOptions[K, V])
			if ok {
				mo.expandable = true
				mo.marshal = nil
				mo.cloneFn = nil
				return mo
			}
		}
	} else if len(ops) == 1 {
		if mo, ok := ops[0].(MapOptions[K, V]); ok {
			if mo.cloneFn != nil && mo.marshal != nil {
				mo.expandable = true
				return mo
			}
		}
	}
	panic(errors.New("invalid map options"))
}
