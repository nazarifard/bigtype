package basic

// import (
// 	"github.com/nazarifard/fastape"
// 	"github.com/nazarifard/syncpool"
// )

// var bufferPool = syncpool.NewBufferPool()

// type MarshalTap[V any, VTape fastape.Tape[V]] struct {
// 	//pool syncpool.BufferPool
// }

// func (t MarshalTap[V, VTape]) Encode(v V) (buf syncpool.Buffer, err error) {
// 	var vTape VTape
// 	size := vTape.Sizeof(v)
// 	buf = bufferPool.Get(size)
// 	//buf.Reset()
// 	_, err = vTape.Roll(v, buf.Bytes())
// 	if err != nil {
// 		buf.Free()
// 	}
// 	return
// }

// func (t MarshalTap[V, VTape]) Decode(bs []byte) (v V, n int, err error) {
// 	var vTape VTape
// 	n, err = vTape.Unroll(bs, &v)
// 	return
// }

// // func NewTap[V any, VTape fastape.Tape[V]]() fastape.Tap[V] {
// // 	return MarshalTap[V, VTape]{}
// // }

// var MTapeString = MarshalTap[string, fastape.StringTape]{}
// var MTapeUint64 = MarshalTap[uint64, fastape.UnitTape[uint64]]{}
