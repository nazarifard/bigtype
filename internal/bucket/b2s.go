package bucket

import (
	"unsafe"
)

// func unsafeCloneBytes(orig []byte) (copy []byte) {
// 	return s2b(b2s(orig))
// }

func b2s(b []byte) string {
	p := unsafe.SliceData(b)
	return unsafe.String(p, len(b))
	//return unsafe.String(&b[0], len(b))
}

func s2b(s string) []byte {
	p := unsafe.StringData(s)
	b := unsafe.Slice(p, len(s))
	return b
}

// //go:noinline
// func hdrEncode(h CellHeader) []byte {
// 	// leaking parameter (closure)
// 	return (*[CellHeaderSize]byte)(unsafe.Pointer(&h))[:]
// }

// func hdrDecode(bs []byte) CellHeader {

//does v have method M ?
// func check(v interface{}) bool {
// 	_, has := v.(interface{M(int) string})
// 	return has
// }

//check N>=M compile time
//var _ = map[bool]struct{}{false: struct{}{}, N>=M: struct{}{}}

//not blank?
//const _ = 1/len(aStringConstant)

//Checks if the generic T is indeed one of the comparable types:
//reflect.TypeOf((*T)(nil)).Elem().Comparable()
