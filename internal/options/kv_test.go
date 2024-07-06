package options

import (
	"math/rand"
	"testing"
	"unsafe"

	marshal "github.com/nazarifard/marshaltap"
	"github.com/nazarifard/marshaltap/tap/stdlib/gob"
)

type IntString struct {
	Key   int
	Value string
}

type IntBytes struct {
	Key   int
	Value []byte
}

type StringInt struct {
	Key   string
	Value int
}

type BytesInt struct {
	Key   []byte
	Value int
}

type StringString struct {
	Key   string
	Value string
}

type StringBytes struct {
	Key   string
	Value []byte
}

func TestKV_IntString(t *testing.T) {
	m := NewKVmarshal[int, string, IntString](nil, nil)
	p := IntString{123456, "qwertyyuuiop"}
	var q IntString
	buf, err := m.Encode(p)
	q, n, err := m.Decode(buf.Bytes())
	if p != q {
		t.Errorf("p!=q %v!=%v", p, q)
	}
	_, _, _ = q, n, err
}

func TestKV_StringInt(t *testing.T) {
	m := NewKVmarshal[string, int, StringInt](nil, nil)
	p := StringInt{"qwertyyuuiop", 123456}
	var q StringInt
	buf, err := m.Encode(p)
	q, n, err := m.Decode(buf.Bytes())
	if p != q {
		t.Errorf("p!=q %v!=%v", p, q)
	}
	_, _, _ = q, n, err
}

func TestKV_IntBytes(t *testing.T) {
	m := NewKVmarshal[int, []byte, IntBytes](nil, nil)
	p := IntBytes{123456, []byte("1234567890")}
	var q IntBytes
	buf, err := m.Encode(p)
	q, n, err := m.Decode(buf.Bytes())
	if p.Key != q.Key || string(p.Value) != string(q.Value) {
		t.Errorf("p!=q %v!=%v", p, q)
	}
	_, _, _ = q, n, err
}

// func TestKV_BytesInt(t *testing.T) {
// 	m := NewKVmarshal[[]byte, int, BytesInt](nil, nil)
// 	p := BytesInt{[]byte("qwertyyuuiop"), 123456}
// 	var q BytesInt
// 	buf, err := m.Encode(p)
// 	q, n, err := m.Decode(buf.Bytes())
// 	if string(p.Key) != string(q.Key) || p.Value != q.Value {
// 		t.Errorf("p!=q %v!=%v", p, q)
// 	}
// 	_, _, _ = q, n, err
// }

func TestKV_StringString(t *testing.T) {
	m := NewKVmarshal[string, string, StringString](nil, nil)
	p := StringString{"qwertyyuuiop", "123456"}
	var q StringString
	buf, err := m.Encode(p)
	q, n, err := m.Decode(buf.Bytes())
	if p.Key != q.Key || p.Value != q.Value {
		t.Errorf("p!=q %v!=%v", p, q)
	}
	_, _, _ = q, n, err
}

func b2s(b []byte) string {
	p := unsafe.SliceData(b)
	return unsafe.String(p, len(b))
	//return unsafe.String(&b[0], len(b))
}
func Benchmark_KV_StringString(b *testing.B) {
	m := NewKVmarshal[string, string, StringString](nil, nil)
	p := StringString{"12345678901234567890", "12345678901234567890"}
	var q StringString
	var n int
	k := []byte(p.Key)
	v := []byte(p.Value)
	for range b.N {
		p.Key = b2s(k[:rand.Int31n(10)])
		p.Value = b2s(v[:rand.Int31n(16)])
		buf, err := m.Encode(p)
		_ = err
		q, n, err = m.Decode(buf.Bytes())
		buf.Free()
		_, _, _ = q, n, err
	}
}

func TestKV_StringBytes(t *testing.T) {
	m := NewKVmarshal[string, []byte, StringBytes](nil, nil)
	p := StringBytes{"qwertyyuuiop", []byte("123456")}
	var q StringBytes
	buf, err := m.Encode(p)
	q, n, err := m.Decode(buf.Bytes())
	if p.Key != q.Key || string(p.Value) != string(q.Value) {
		t.Errorf("p!=q %v!=%v", p, q)
	}
	_, _, _ = q, n, err
}

type Point struct {
	X            float32
	Y            float32
	LocationName string
}

// var tk marshal.Interface[K] = gob.GobTap[K]{}
// var tv marshal.Interface[V] = gob.GobTap[V]{}
func TestKV_StringPoint(t *testing.T) {
	var tv marshal.Interface[Point] = gob.GobTap[Point]{}
	m := NewKVmarshal[string, Point, struct {
		Key   string
		Value Point
	}](nil, tv)
	p := struct {
		Key   string
		Value Point
	}{Key: "123456", Value: Point{1.11, 2.22, "abcd"}}
	var q struct {
		Key   string
		Value Point
	}
	buf, err := m.Encode(p)
	if err != nil {
		t.Errorf("encode Failed")
	}
	q, n, err := m.Decode(buf.Bytes())
	if p != q {
		t.Errorf("p!=q %v!=%v", p, q)
	}
	_, _, _ = q, n, err
}

// var tk marshal.Interface[K] = gob.GobTap[K]{}
// var tv marshal.Interface[V] = gob.GobTap[V]{}
func TestKV_PointString(t *testing.T) {
	var tk marshal.Interface[Point] = gob.GobTap[Point]{}
	m := NewKVmarshal[Point, string, struct {
		Key   Point
		Value string
	}](tk, nil)
	p := struct {
		Key   Point
		Value string
	}{Value: "123456", Key: Point{1.11, 2.22, "abcd"}}
	var q struct {
		Key   Point
		Value string
	}
	buf, err := m.Encode(p)
	if err != nil {
		t.Errorf("encode Failed")
	}
	q, n, err := m.Decode(buf.Bytes())
	if p != q {
		t.Errorf("p!=q %v!=%v", p, q)
	}
	_, _, _ = q, n, err
}

func TestKV_PointPoint(t *testing.T) {
	var tk marshal.Interface[Point] = gob.GobTap[Point]{}
	var tv marshal.Interface[Point] = gob.GobTap[Point]{}
	m := NewKVmarshal[Point, Point, struct {
		Key   Point
		Value Point
	}](tk, tv)
	p := struct {
		Key   Point
		Value Point
	}{Value: Point{1.111, 2.222, "xyz"}, Key: Point{1.11, 2.22, "abcd"}}
	var q struct {
		Key   Point
		Value Point
	}
	buf, err := m.Encode(p)
	if err != nil {
		t.Errorf("encode Failed")
	}
	q, n, err := m.Decode(buf.Bytes())
	if p != q {
		t.Errorf("p!=q %v!=%v", p, q)
	}
	_, _, _ = q, n, err
}
