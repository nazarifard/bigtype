package basic

import (
	"errors"

	marshal "github.com/nazarifard/marshaltap"
)

type array[V any] struct {
	*BytesArray
	marshal.Interface[V]
	cloneFn func(V) V
}

var ErrInvalidData = errors.New("invalid data")

// must be released by marshal.Interface after used
func (ba array[V]) UnsafePtr(index int) *V {
	pbs := ba.BytesArray.UnsafePtr(index)
	if pbs == nil {
		return nil
	}
	v, _, err := ba.Interface.Decode(*pbs)
	if err != nil {
		panic(ErrInvalidData)
	}
	return v
}

func (ba array[V]) Get(index int) (v V) {
	pv := ba.UnsafePtr(index)
	if pv == nil {
		return
	}
	defer ba.Interface.Free(pv)
	return ba.cloneFn(*pv)
}

func (ba array[V]) Set(index int, v V) {
	buf, err := ba.Interface.Encode(v)
	defer buf.Free()
	if err != nil {
		panic(ErrInvalidData)
	}
	ba.BytesArray.Set(index, buf.Bytes())
}

func (ba *array[V]) Update(index int, fn func(V) V) {
	old := Val(ba.BytesArray.UnsafePtr(index))
	v, _, err := ba.Interface.Decode(old)
	defer ba.Interface.Free(v)
	if err != nil {
		panic(ErrInvalidData)
	}

	buf, err := ba.Interface.Encode(fn(*v)) //here we dont need to cloneFn
	defer buf.Free()
	if err != nil {
		panic(ErrInvalidData)
	}
	ba.BytesArray.Set(index, buf.Bytes())
}

func NewArray[V any](ops ...any) Array[V] {
	ao := ParsArrayOptions[V](ops...)
	var v V
	if IsFixedType(v) {
		return NewFixedArray[V](ao.Size, ao.expandable)
	}

	if IsBytes(v) {
		return any(NewBytesArray(ao.Size, ao.expandable)).(Array[V])
	}

	if IsString(v) {
		return any(NewStringArray(ao.Size, ao.expandable)).(Array[V])
	}

	return &array[V]{
		BytesArray: NewBytesArray(ao.Size, ao.expandable),
		cloneFn:    ao.cloneFn,
		Interface:  ao.marshal,
	}
}
