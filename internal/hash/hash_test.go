package hash

import (
	"testing"
	"unsafe"
)

func Benchmark_HashArray(b *testing.B) {
	var v [12]byte
	copy(v[:], "123456789012")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		p := (*byte)(unsafe.Pointer(&v))
		bs := unsafe.Slice(p, unsafe.Sizeof(v))

		if len(bs) != 12 {
			b.Errorf("got:%s want:%s", bs, v)
		}
	}
}

// func BenchmarkHash_ByteSlice_T1ha0(b *testing.B) {
// 	hash := HashBuilder[[]byte]()
// 	k := []byte("123456789012345678901234567890")
// 	var g uint64
// 	for i := 0; i < b.N; i++ {
// 		h := hash.Hash(k)
// 		if i == 0 {
// 			g = h
// 		}
// 		if g != h {
// 			b.Errorf("g!=h %d!=%d", g, h)
// 		}
// 		g = h
// 	}
// }
