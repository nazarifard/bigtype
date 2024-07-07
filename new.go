package bigtype

import "github.com/nazarifard/bigtype/internal/basic"

func NewArray[V any](ops ...any) Array[V] {
	return basic.NewArray[V](ops...)
}

func NewMap[K comparable, V any](ops ...any) Map[K, V] {
	return basic.NewMap[K, V](ops...)
}
