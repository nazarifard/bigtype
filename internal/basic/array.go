package basic

import (
	"errors"

	"github.com/nazarifard/fastape"
)

type array[V any] struct {
	*BytesArray
	fastape.Tape[V]
}

var ErrInvalidData = errors.New("invalid data")

func (ba array[V]) Get(index int) (v V) {
	bs := ba.BytesArray.UnsafePtr(index)
	ba.Tape.Unroll(bs, &v)
	return
}

func (ba array[V]) Set(index int, v V) {
	space := ba.BytesArray.Request(index, ba.Tape.Sizeof(v))
	ba.Tape.Roll(v, space)
}

func (ba *array[V]) Update(index int, fn func(V) V) {
	var old V
	bs := ba.BytesArray.UnsafePtr(index)
	_, err := ba.Tape.Unroll(bs, &old)
	if err != nil {
		panic(ErrInvalidData)
	}
	now := fn(old)
	//defer ba.Tape.Free(old)
	ba.Set(index, now)
	//ba.Tape.Free(now)
}

func NewArray[V any](ops ...any) Array[V] {
	ao := ParsArrayOptions[V](ops...)
	var v V
	if IsFixedType(v) {
		return NewFixedArray[V](ao.Size(), ao.Expandable())
	}

	if IsBytes(v) {
		return any(NewBytesArray(ao.Size(), ao.Expandable())).(Array[V])
	}

	if IsString(v) {
		return any(NewStringArray(ao.Size(), ao.Expandable())).(Array[V])
	}

	if ao.Marshal() == nil {
		//TODO
		//ao.Tape = fastape.GobTape
		panic("marshal tape is nil")
	}
	return &array[V]{
		BytesArray: NewBytesArray(ao.Size(), ao.Expandable()),
		Tape:       ao.Marshal(),
	}
}
