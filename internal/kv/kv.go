package kv

import (
	"unsafe"

	"github.com/nazarifard/bigtype/internal/utils"
	"github.com/nazarifard/fastape"
)

type Type int

const (
	Invalid Type = iota
	Bytes
	String
	Fixed
	Unknown
)

type KV[K comparable, V any] struct {
	Key   K
	Value V
}

type TapeKV[K comparable, V any] struct {
	vTape fastape.Tape[V]
	kTape fastape.Tape[K]
	vType Type
	kType Type
}

func (t TapeKV[K, V]) Roll(kv KV[K, V], bs []byte) (n int, err error) {
	if t.kType == Fixed {
		//k is fixed
		pk := (*byte)(unsafe.Pointer(&kv.Key))
		n += copy(bs, unsafe.Slice(pk, int(unsafe.Sizeof(kv.Key))))
		m, err := t.VRoll(kv.Value, bs[n:])
		if err != nil {
			return 0, err
		}
		n += m
		return n, err
	}

	if t.vType == Fixed {
		//k is fixed
		pk := (*byte)(unsafe.Pointer(&kv.Value))
		n += copy(bs, unsafe.Slice(pk, int(unsafe.Sizeof(kv.Value))))

		m, err := t.KRoll(kv.Key, bs[n:])
		if err != nil {
			return 0, err
		}
		n += m
		return n, err
	}

	n, err = t.vTape.Roll(kv.Value, bs)
	if err != nil {
		return 0, err
	}

	switch t.kType {
	case String:
		s := (*string)(unsafe.Pointer(&kv.Key))
		n += copy(bs[n:], *s)
	case Unknown:
		m, err := t.kTape.Roll(kv.Key, bs[n:])
		if err != nil {
			return 0, err
		}
		n += m
	}

	return
}

func (t TapeKV[K, V]) KRoll(k K, bs []byte) (n int, err error) {
	switch t.kType {
	case String:
		s := (*string)(unsafe.Pointer(&k))
		n += copy(bs, *s)
	case Unknown:
		m, err := t.kTape.Roll(k, bs)
		if err != nil {
			return 0, err
		}
		n += m
	}
	return
}

func (t TapeKV[K, V]) VRoll(v V, bs []byte) (n int, err error) {
	switch t.vType {
	case Bytes:
		b := *(*[]byte)(unsafe.Pointer(&v))
		n += copy(bs, b)
	case String:
		s := (*string)(unsafe.Pointer(&v))
		n += copy(bs, *s)
	case Unknown:
		m, err := t.vTape.Roll(v, bs)
		if err != nil {
			return 0, err
		}
		n += m
	}
	return
}

func (t TapeKV[K, V]) VUnrol(bs []byte, p *V) (n int, err error) {
	switch t.vType {
	case Bytes:
		b := *(*[]byte)(unsafe.Pointer(p))
		n += copy(b, bs)
	case String:
		s := (*string)(unsafe.Pointer(p))
		*s = string(bs)
		n += len(*s)
	case Unknown:
		m, err := t.vTape.Unroll(bs, p)
		if err != nil {
			return 0, err
		}
		n += m
	}
	return
}

func (t TapeKV[K, V]) KUnrol(bs []byte, p *K) (n int, err error) {
	switch t.kType {
	case String:
		s := (*string)(unsafe.Pointer(p))
		*s = string(bs)
		n += len(*s)
	case Unknown:
		m, err := t.kTape.Unroll(bs, p)
		if err != nil {
			return 0, err
		}
		n += m
	}
	return
}

func (t TapeKV[K, V]) Unroll(bs []byte, pkv *KV[K, V]) (n int, err error) {
	if t.vType == Fixed {
		p := (*byte)(unsafe.Pointer(&pkv.Value))
		n += copy(unsafe.Slice(p, int(unsafe.Sizeof(pkv.Value))), bs)

		m, err := t.KUnrol(bs[n:], &pkv.Key)
		if err != nil {
			return 0, err
		}
		n += m
		return n, err
	}

	if t.kType == Fixed {
		//k is fixed
		pk := (*byte)(unsafe.Pointer(&pkv.Key))
		n += copy(unsafe.Slice(pk, int(unsafe.Sizeof(pkv.Key))), bs)

		m, err := t.VUnrol(bs[n:], &pkv.Value)
		if err != nil {
			return 0, err
		}
		n += m
		return n, err
	}

	n, err = t.vTape.Unroll(bs[n:], &pkv.Value)
	if err != nil {
		return 0, err
	}

	switch t.kType {
	case String:
		s := (*string)(unsafe.Pointer(&pkv.Key))
		*s = string(bs[n:])
		n += len(*s)
	case Unknown:
		m, err := t.kTape.Unroll(bs[n:], &pkv.Key)
		if err != nil {
			return 0, err
		}
		n += m
	}

	return n, err
}

func (t TapeKV[K, V]) Sizeof(kv KV[K, V]) int {
	var n int
	if t.kType == Fixed {
		//k is fixed
		n = int(unsafe.Sizeof(kv.Key))
		switch t.vType {
		case Bytes:
			bs := *(*[]byte)(unsafe.Pointer(&kv.Value))
			n += len(bs)
		case String:
			s := (*string)(unsafe.Pointer(&kv.Value))
			n += len(*s)
		case Unknown:
			n += t.vTape.Sizeof(kv.Value)
		}
		return n
	}

	if t.vType == Fixed {
		//k is fixed
		n = int(unsafe.Sizeof(kv.Value))
		switch t.kType {
		case String:
			s := (*string)(unsafe.Pointer(&kv.Key))
			n += len(*s)
		case Unknown:
			n += t.kTape.Sizeof(kv.Key)
		}
		return n
	}

	n = t.vTape.Sizeof(kv.Value)
	switch t.kType {
	case String:
		s := (*string)(unsafe.Pointer(&kv.Key))
		n += len(*s)
	case Unknown:
		n += t.kTape.Sizeof(kv.Key)
	}
	return n
}

func NewTapeKV[K comparable, V any](kTape fastape.Tape[K], vTape fastape.Tape[V]) TapeKV[K, V] {
	tape := TapeKV[K, V]{
		vTape: vTape,
		kTape: kTape,
		vType: getType(*new(V)),
		kType: getType(*new(K)),
	}
	switch tape.kType {
	case Fixed:
	case String:
		tape.kTape = any(fastape.StringTape{}).(fastape.Tape[K])
	case Unknown:
	}

	switch tape.vType {
	case Fixed:
	case Bytes:
		tape.vTape = any(fastape.SliceTape[byte, fastape.UnitTape[byte]]{}).(fastape.Tape[V])
	case String:
		tape.vTape = any(fastape.StringTape{}).(fastape.Tape[V])
	case Unknown:
	}
	return tape
}

func getType[V any](v V) Type {
	var tipe Type
	if utils.IsBytes(v) {
		tipe = Bytes
	} else if utils.IsString(v) {
		tipe = String
	} else if utils.IsFixedType(v) {
		tipe = Fixed
	}
	return tipe
}

func (t TapeKV[K, V]) Cat(k K, v V, bs []byte) (n int, err error) {
	if t.kType == Fixed {
		//k is fixed
		pk := (*byte)(unsafe.Pointer(&k))
		n += copy(bs, unsafe.Slice(pk, int(unsafe.Sizeof(k))))
		m, err := t.VRoll(v, bs[n:])
		if err != nil {
			return 0, err
		}
		n += m
		return n, err
	}

	if t.vType == Fixed {
		//k is fixed
		pk := (*byte)(unsafe.Pointer(&v))
		n += copy(bs, unsafe.Slice(pk, int(unsafe.Sizeof(v))))

		m, err := t.KRoll(k, bs[n:])
		if err != nil {
			return 0, err
		}
		n += m
		return n, err
	}

	n, err = t.vTape.Roll(v, bs)
	if err != nil {
		return 0, err
	}

	switch t.kType {
	case String:
		s := (*string)(unsafe.Pointer(&k))
		n += copy(bs[n:], *s)
	case Unknown:
		m, err := t.kTape.Roll(k, bs[n:])
		if err != nil {
			return 0, err
		}
		n += m
	}

	return
}

func (t TapeKV[K, V]) SizeOf(k K, v V) int {
	var n int
	if t.kType == Fixed {
		//k is fixed
		n = int(unsafe.Sizeof(k))
		switch t.vType {
		case Bytes:
			bs := *(*[]byte)(unsafe.Pointer(&v))
			n += len(bs)
		case String:
			s := (*string)(unsafe.Pointer(&v))
			n += len(*s)
		case Unknown:
			n += t.vTape.Sizeof(v)
		}
		return n
	}

	if t.vType == Fixed {
		//k is fixed
		n = int(unsafe.Sizeof(v))
		switch t.kType {
		case String:
			s := (*string)(unsafe.Pointer(&k))
			n += len(*s)
		case Unknown:
			n += t.kTape.Sizeof(k)
		}
		return n
	}

	n = t.vTape.Sizeof(v)
	switch t.kType {
	case String:
		s := (*string)(unsafe.Pointer(&k))
		n += len(*s)
	case Unknown:
		n += t.kTape.Sizeof(k)
	}
	return n
}

func (t TapeKV[K, V]) KeyValue(bs []byte) (key K, value V, ok bool) {
	value, ok1 := t.Value(bs)
	key, ok2 := t.Key(bs)
	return key, value, ok1 && ok2
}

func (t TapeKV[K, V]) Value(bs []byte) (V, bool) {
	var v V
	//var ok bool
	if t.vType == Fixed {
		p := (*byte)(unsafe.Pointer(&v))
		copy(unsafe.Slice(p, int(unsafe.Sizeof(v))), bs)
		return v, true
	}

	if t.kType == Fixed {
		//k is fixed
		var v2 V
		n := int(unsafe.Sizeof(*new(K)))
		_, err := t.VUnrol(bs[n:], &v2)
		if err != nil {
			return *new(V), false
		}
		v = v2
		return v, true
	}

	switch t.vType {
	case Bytes:
		var lenTape fastape.LenTape
		var size int
		n, err := lenTape.Unroll(bs, &size)
		if err != nil {
			return v, false
		}

		return any(bs[n : n+size]).(V), true
	case String:
		var lenTape fastape.LenTape
		var size int
		n, err := lenTape.Unroll(bs, &size)
		if err != nil {
			return v, false
		}
		return any(unsafe.String(&bs[n], size)).(V), true
	case Unknown:
		return func() (V, bool) {
			var v2 V
			_, err := t.vTape.Unroll(bs, &v2)
			if err != nil {
				return v2, false
			}
			return v2, true
		}()
	}
	return v, false
}

func (t TapeKV[K, V]) Key(bs []byte) (k K, ok bool) {
	if t.kType == Fixed {
		p := (*byte)(unsafe.Pointer(&k))
		copy(unsafe.Slice(p, int(unsafe.Sizeof(k))), bs)
		return k, true
	}

	if t.vType == Fixed {
		//k is fixed
		n := int(unsafe.Sizeof(*new(V)))
		_, err := t.KUnrol(bs[n:], &k)
		if err != nil {
			return *new(K), false
		}
		return k, true
	}

	var lenTape fastape.LenTape
	var size int
	n, err := lenTape.Unroll(bs, &size)
	if err != nil {
		return *new(K), false
	}

	switch t.kType {
	case Unknown:
		_, err := t.kTape.Unroll(bs[n+size:], &k)
		if err != nil {
			return *new(K), false
		}
		return k, true
	case Bytes:
		return any(bs[n+size:]).(K), true
	case String:
		return any(unsafe.String(&bs[n+size], len(bs)-(n+size))).(K), true
	}

	return k, false
}
