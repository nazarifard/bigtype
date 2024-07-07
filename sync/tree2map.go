package sync

import (
	"fmt"
	"reflect"
)

func isNumber[V any]() bool {
	switch reflect.ValueOf(*new(V)).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

func makeTree[K comparable, V any](ops ...any) Map[K, V] {
	switch reflect.ValueOf(*new(K)).Kind() {
	case reflect.Int:
		t := newTree[int, V](ops...)
		return any(t).(Map[K, V])
	case reflect.Int8:
		t := newTree[int8, V](ops...)
		return any(t).(Map[K, V])
	case reflect.Int16:
		t := newTree[int16, V](ops...)
		return any(t).(Map[K, V])
	case reflect.Int32:
		t := newTree[int32, V](ops...)
		return any(t).(Map[K, V])
	case reflect.Int64:
		t := newTree[int64, V](ops...)
		return any(t).(Map[K, V])
	case reflect.Uint:
		t := newTree[uint, V](ops...)
		return any(t).(Map[K, V])
	case reflect.Uint8:
		t := newTree[uint8, V](ops...)
		return any(t).(Map[K, V])
	case reflect.Uint16:
		t := newTree[uint16, V](ops...)
		return any(t).(Map[K, V])
	case reflect.Uint32:
		t := newTree[uint32, V](ops...)
		return any(t).(Map[K, V])
	case reflect.Uint64:
		t := newTree[uint64, V](ops...)
		return any(t).(Map[K, V])
	case reflect.Float32:
		t := newTree[float32, V](ops...)
		return any(t).(Map[K, V])
	case reflect.Float64:
		t := newTree[float64, V](ops...)
		return any(t).(Map[K, V])
	default:
		panic(fmt.Errorf("MakeTree() failed. K is not comparable. Type: %T", *new(K)))
	}
}
