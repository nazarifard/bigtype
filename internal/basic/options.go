package basic

import (
	"errors"

	"github.com/nazarifard/bigtype/options"
)

func ParsArrayOptions[V any](ops ...any) (ao options.ArrayOptions[V]) {
	if len(ops) == 1 {
		switch t := ops[0].(type) {
		case int:
			ao.WithSize(t).WithExpandable(false)
		case options.ArrayOptions[V]:
			ao = t
		case *options.ArrayOptions[V]:
			ao = *t
		default:
			panic(errors.New("invalid array options"))
		}
	}
	return ao
}

func ParsTreeOptions[K kNumber, V any](ops ...any) (to options.TreeOptions[K, V]) {
	if len(ops) == 1 {
		switch t := ops[0].(type) {
		case int:
			to.WithSize(t)
		case options.TreeOptions[K, V]:
			to = t
		case *options.TreeOptions[K, V]:
			to = *t
		default:
			panic(errors.New("invalid array options"))
		}
	}
	return to
}

func ParsMapOptions[K comparable, V any](ops ...any) (mo options.MapOptions[K, V]) {
	var v V
	if !IsRequiredToMarshal(v) {
		switch len(ops) {
		case 0:
			return mo
		case 1:
			size, ok := ops[0].(int)
			if ok {
				mo.WithSize(size)
				return mo
			}

			mo, ok := ops[0].(options.MapOptions[K, V])
			if ok {
				mo.WithMarshal(nil)
				return mo
			}
			if mo, ok := ops[0].(*options.MapOptions[K, V]); ok {
				mo.WithMarshal(nil)
				return *mo
			}
		}
	} else if len(ops) == 1 {
		if mo, ok := ops[0].(options.MapOptions[K, V]); ok {
			if mo.Marshal() != nil {
				return mo
			}
		}
		if mo, ok := ops[0].(*options.MapOptions[K, V]); ok {
			if mo.Marshal() != nil {
				return *mo
			}
		}
	}
	panic(errors.New("invalid map options"))
}
