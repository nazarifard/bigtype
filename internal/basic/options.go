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
	switch len(ops) {
	case 0:
	case 1:
		if size, ok := ops[0].(int); ok {
			mo.WithSize(size)
		} else {
			var ok bool
			if mo, ok = ops[0].(options.MapOptions[K, V]); !ok {
				if pmo, ok := ops[0].(*options.MapOptions[K, V]); ok {
					mo = *pmo
				}
			}
		}
	}

	if IsRequiredToMarshal(*new(V)) {
		if mo.VTape() == nil {
			panic(errors.New("invalid map options. VTape is missed"))
		}
	}

	if IsRequiredToMarshal(*new(K)) {
		if mo.KTape() == nil {
			panic(errors.New("invalid map options. KTape is missed"))
		}
	}

	return
}
