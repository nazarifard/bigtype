package hash

import (
	"unsafe"

	"github.com/nazarifard/bigtype/internal/utils"
	tap "github.com/nazarifard/marshaltap"
)

const defaultSeed uint64 = 0x7d5b016bcbfebb4c

type Hashable[K any] interface {
	Hash(k K) uint64
}

type StringHash struct {
	seed uint64
}

func (h *StringHash) Hash(s string) uint64 {
	p := unsafe.StringData(s)
	bs := unsafe.Slice(p, len(s))
	return T1ha0(bs, h.seed)
}

type BytesHash struct {
	seed uint64
}

func (h *BytesHash) Hash(bs []byte) uint64 {
	return T1ha0(bs, h.seed)
}

type FixObject[K any] struct {
	seed uint64
}

func (h *FixObject[K]) Hash(k K) uint64 {
	p := (*byte)(unsafe.Pointer(&k))
	bs := unsafe.Slice(p, unsafe.Sizeof(*new(K)))
	return T1ha0(bs, h.seed)
}

type Hash[K any] struct {
	seed   uint64
	hashFn func(K, uint64) uint64
}

func (h *Hash[K]) Hash(k K) uint64 {
	return h.hashFn(k, h.seed)
}

type HashM[K any, M tap.Tap[K]] struct {
	seed      uint64
	Marshaler M
}

func (h *HashM[K, M]) Hash(k K) uint64 {
	buf, _ := h.Marshaler.Encode(k)
	n64 := T1ha0(buf.Bytes(), h.seed)
	buf.Free()
	return n64
}

func NewHash[K any](m ...tap.Tap[K]) Hashable[K] {
	var k K
	if utils.IsBytes(k) {
		h := &BytesHash{
			seed: defaultSeed,
		}
		return any(h).(Hashable[K])

	} else if utils.IsString(k) {
		h := &StringHash{
			seed: defaultSeed,
		}
		return any(h).(Hashable[K])

	} else if utils.IsFixedType(k) {
		h := &FixObject[K]{seed: defaultSeed}
		return any(h).(Hashable[K])

	} else if len(m) == 1 {
		h := &HashM[K, tap.Tap[K]]{
			seed: defaultSeed,
		}
		return any(h).(Hashable[K])
	} else {
		var t tap.Tap[K] = nil //gob.GobTap[K]{}
		h := &HashM[K, tap.Tap[K]]{
			seed:      defaultSeed,
			Marshaler: t,
		}
		return any(h).(Hashable[K])
	}
}
