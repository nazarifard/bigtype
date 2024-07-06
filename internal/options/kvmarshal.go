package options

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/nazarifard/fastape"
	marshal "github.com/nazarifard/marshaltap"
	"github.com/nazarifard/syncpool"
)

type elementKind int

const (
	Unknown = iota
	FixedSize
	String
	Bytes
)

type KVMarshal[K comparable, V any, KV ~struct {
	Key   K
	Value V
}] struct {
	tk marshal.Interface[K]
	tv marshal.Interface[V]
	//kTape   *fastape.UnitTape[K]
	//vTape   *fastape.UnitTape[V]
	vKind      elementKind
	kKind      elementKind
	bufferPool syncpool.BufferPool
}

func NewKVmarshal[K comparable, V any, KV ~struct {
	Key   K
	Value V
}](tk marshal.Interface[K], tv marshal.Interface[V]) *KVMarshal[K, V, KV] {
	kv := KVMarshal[K, V, KV]{}
	kv.tk = tk
	kv.tv = tv
	kv.bufferPool = syncpool.NewBufferPool()

	if IsFixedType(*new(K)) {
		kv.kKind = FixedSize
	} else {
		switch reflect.ValueOf(*new(K)).Kind() {
		case reflect.String:
			kv.kKind = String
		case reflect.Slice:
			switch reflect.TypeOf(*new(K)).Elem().Kind() {
			case reflect.Uint8: //, reflect.Int8, reflect.Bool
				kv.kKind = Bytes
			}
		default:
			kv.kKind = Unknown
		}
	}

	if IsFixedType(*new(V)) {
		kv.vKind = FixedSize
	} else {
		switch reflect.ValueOf(*new(V)).Kind() {
		case reflect.String:
			kv.vKind = String
		case reflect.Slice:
			switch reflect.TypeOf(*new(V)).Elem().Kind() {
			case reflect.Uint8: //, reflect.Int8, reflect.Bool
				kv.vKind = Bytes
			}
		default:
			kv.vKind = Unknown
		}
	}

	if kv.vKind == FixedSize && kv.kKind == FixedSize {
		fmt.Println("warning: both key and value are fixed type. ")
		fmt.Println("for fixed type shoudnt use marshal")
		return &kv
	}
	return &kv
}

func s2b(s string) []byte {
	p := unsafe.StringData(s)
	b := unsafe.Slice(p, len(s))
	return b
}

func (t KVMarshal[K, V, KV]) getK(k K, kind elementKind) []byte {
	switch kind {
	case FixedSize:
		b := (*byte)(unsafe.Pointer(&k))
		return unsafe.Slice(b, unsafe.Sizeof(*new(K)))
	case Bytes:
		bs := *(*[]byte)(unsafe.Pointer(&k))
		return bs
	case String:
		s := (*string)(unsafe.Pointer(&k))
		return s2b(*s)
	default:
		panic("unexpected type recieved KVMarshal.EncodeElement() function")
	}
}

func (t KVMarshal[K, V, KV]) getV(v V, kind elementKind) []byte {
	switch kind {
	case FixedSize:
		b := (*byte)(unsafe.Pointer(&v))
		return unsafe.Slice(b, unsafe.Sizeof(*new(K)))
	case Bytes:
		bs := *(*[]byte)(unsafe.Pointer(&v))
		return bs
	case String:
		s := (*string)(unsafe.Pointer(&v))
		return s2b(*s)
	default:
		panic("unexpected type recieved KVMarshal.EncodeElement() function")
	}
}

func (t KVMarshal[K, V, KV]) Encode(kv KV) (buf syncpool.Buffer, err error) {
	k := struct {
		Key   K
		Value V
	}(kv).Key
	v := struct {
		Key   K
		Value V
	}(kv).Value

	if t.kKind == Unknown && t.vKind == Unknown {
		buf, err = t.tk.Encode(k)
		if err != nil {
			buf.Free()
			return
		}
		vbuf, err2 := t.tv.Encode(v)
		if err2 != nil {
			buf.Free()
			vbuf.Free()
			return nil, err2
		}
		buf.Write(vbuf.Bytes())
		vbuf.Free()
		return
	} else if t.kKind == Unknown && t.vKind != Unknown {
		buf, err = t.tk.Encode(k)
		if err != nil {
			buf.Free()
			return
		}
		buf.Write(t.getV(v, t.vKind))
		return
	} else if t.kKind != Unknown && t.vKind == Unknown {
		buf, err = t.tv.Encode(v)
		if err != nil {
			buf.Free()
			return
		}
		buf.Write(t.getK(k, t.kKind))
		return
	}

	if t.kKind == FixedSize {
		kTape := fastape.UnitTape[K]{}
		buf = t.bufferPool.Get(int(unsafe.Sizeof(k)))
		kTape.Roll(k, buf.Bytes())
		//append v
		switch t.vKind {
		case Bytes:
			bs := *(*[]byte)(unsafe.Pointer(&v))
			buf.Write(bs)
			return
		case String:
			s := (*string)(unsafe.Pointer(&v))
			buf.WriteString(*s)
			return
		}
	} else if t.vKind == FixedSize {
		vTape := fastape.UnitTape[V]{}
		buf = t.bufferPool.Get(int(unsafe.Sizeof(v)))
		vTape.Roll(v, buf.Bytes())
		//append k
		switch t.kKind {
		case Bytes:
			bs := *(*[]byte)(unsafe.Pointer(&k))
			buf.Write(bs)
			return
		case String:
			s := (*string)(unsafe.Pointer(&k))
			buf.WriteString(*s)
			return
		}
	}
	//{ //both k,v dynamic size
	var strTape fastape.StringTape
	var bsTape fastape.SliceTape[byte, fastape.UnitTape[byte]]
	switch t.kKind {
	case Bytes:
		bs := *(*[]byte)(unsafe.Pointer(&k))
		kSize := bsTape.Sizeof(bs)
		buf = t.bufferPool.Get(kSize)
		bsTape.Roll(bs, buf.Bytes())
	case String:
		s := (*string)(unsafe.Pointer(&k))
		kSize := strTape.Sizeof(*s)
		buf = t.bufferPool.Get(kSize)
		strTape.Roll(*s, buf.Bytes())
	}

	switch t.vKind {
	case Bytes:
		bs := *(*[]byte)(unsafe.Pointer(&v))
		buf.Write(bs)
		return
	case String:
		s := (*string)(unsafe.Pointer(&v))
		buf.WriteString(*s)
		return
	}
	return
}

func (t KVMarshal[K, V, KV]) Decode(bs []byte) (kv KV, n int, err error) {
	var k K
	var v V
	if t.kKind == Unknown && t.vKind == Unknown {
		return t.Decode2(bs)

	} else if t.kKind == Unknown && t.vKind != Unknown {
		k, n, err = t.tk.Decode(bs)
		if err != nil {
			return
		}
		bs = bs[n:]
		switch t.vKind {
		case FixedSize:
			v = *(*V)(unsafe.Pointer(&bs[0]))
		case Bytes:
			vbs := (*[]byte)(unsafe.Pointer(&v))
			*vbs = append(*vbs, bs...)
		case String:
			s := (*string)(unsafe.Pointer(&v))
			*s = string(bs)
		}
		return struct {
			Key   K
			Value V
		}{Key: k, Value: v}, n + len(bs), nil

	} else if t.kKind != Unknown && t.vKind == Unknown {
		v, n, err = t.tv.Decode(bs)
		if err != nil {
			return
		}
		bs = bs[n:]
		switch t.kKind {
		case FixedSize:
			k = *(*K)(unsafe.Pointer(&bs[0]))
		case Bytes:
			kbs := (*[]byte)(unsafe.Pointer(&k))
			*kbs = append(*kbs, bs...)
		case String:
			s := (*string)(unsafe.Pointer(&k))
			*s = string(bs)
		}
		return struct {
			Key   K
			Value V
		}{Key: k, Value: v}, n + len(bs), nil
	}

	if t.kKind == FixedSize {
		kTape := fastape.UnitTape[K]{}
		kTape.Unroll(bs, &k)
		n += int(unsafe.Sizeof(k))
		bs := bs[unsafe.Sizeof(k):]
		//append v
		switch t.vKind {
		case Bytes:
			vbs := (*[]byte)(unsafe.Pointer(&v))
			*vbs = append(*vbs, bs...) //copy
			return struct {
				Key   K
				Value V
			}{Key: k, Value: v}, n + len(bs), nil
		case String:
			s := (*string)(unsafe.Pointer(&v))
			*s = string(bs) //copy
			return struct {
				Key   K
				Value V
			}{Key: k, Value: v}, n + len(bs), nil
		}
	} else if t.vKind == FixedSize {
		vTape := fastape.UnitTape[V]{}
		vTape.Unroll(bs, &v)
		n += int(unsafe.Sizeof(v))
		bs := bs[unsafe.Sizeof(v):]
		//append v
		switch t.kKind {
		case Bytes:
			kbs := (*[]byte)(unsafe.Pointer(&k))
			*kbs = append(*kbs, bs...) //copy
			return struct {
				Key   K
				Value V
			}{Key: k, Value: v}, n + len(bs), nil
		case String:
			s := (*string)(unsafe.Pointer(&k))
			*s = string(bs) //copy
			return struct {
				Key   K
				Value V
			}{Key: k, Value: v}, n + len(bs), nil
		}
	}

	//{ //both k,v dynamic size
	var bsTape fastape.SliceTape[byte, fastape.UnitTape[byte]]
	var strTape fastape.StringTape
	switch t.kKind {
	case Bytes:
		kbs := *(*[]byte)(unsafe.Pointer(&k))
		n, err = bsTape.Unroll(bs, &kbs)
		if err != nil {
			return
		}
		bs = bs[n:]
	case String:
		pstr := (*string)(unsafe.Pointer(&k))
		n, err = strTape.Unroll(bs, pstr)
		if err != nil {
			return
		}
		bs = bs[n:]
	}

	switch t.vKind {
	case Bytes:
		vbs := (*[]byte)(unsafe.Pointer(&v))
		*vbs = append(*vbs, bs...)
	case String:
		pstr := (*string)(unsafe.Pointer(&v))
		*pstr = string(bs)
	}
	return struct {
		Key   K
		Value V
	}{Key: k, Value: v}, n + len(bs), nil

}

func (t KVMarshal[K, V, KV]) Decode2(bs []byte) (i KV, n int, err error) {
	var n1, n2 int
	k, n1, err := t.tk.Decode(bs)
	if err != nil {
		return
	}
	v, n2, err := t.tv.Decode(bs[n1:])
	return KV(struct {
		Key   K
		Value V
	}{k, v}), n1 + n2, err
}
