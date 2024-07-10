package basic

import "unsafe"

func s2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func b2s(bs []byte) string {
	return unsafe.String(&bs[0], len(bs))
	//return u.String(u.SliceData(bs), len(bs))
}
