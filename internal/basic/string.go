package basic

import "unsafe"

func s2b(s string) []byte {
	p := unsafe.StringData(s)
	b := unsafe.Slice(p, len(s))
	return b
}

// func b2s(b []byte) string {
// 	p := unsafe.SliceData(b)
// 	return unsafe.String(p, len(b))
// 	//return unsafe.String(&b[0], len(b))
// }
