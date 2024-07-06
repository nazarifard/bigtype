package basic

import (
	"testing"
)

//var _ = func() int { config.Load(); return 0 }()

var id = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2'}
var sample = Sample{0, id, id, id, id}

const fixed_array_max_size int = 7 * 1000 * 1000

//var ba Array[Sample]

func setupFixedArray() *BigFixedArray[Sample] {
	ba := NewFixedArray[Sample](fixed_array_max_size)
	for n := 0; n < fixed_array_max_size; n++ {
		//n := i % uint32(fixed_array_max_size)
		sample.Num = uint32(n)
		ba.Set(n, sample)
	}
	return ba.(*BigFixedArray[Sample])
}

// func TestFixedArray(t *testing.T) {
// 	sample.Num = 123456
// 	ba.Set(123456, sample)
// 	got := ba.Get(123456).Num
// 	want := uint32(123456)
// 	if got != want {
// 		t.Errorf("got %v, want %v", got, want)
// 	}
// }

func BenchmarkInsert(b *testing.B) {
	ba := NewFixedArray[Sample](fixed_array_max_size)
	for i := 0; i < b.N; i++ {
		n := i % fixed_array_max_size
		sample.Num = uint32(n)
		ba.Set(n, sample)
	}
}

func BenchmarkRead(b *testing.B) {
	ba := setupFixedArray()
	for i := 0; i < b.N; i++ {
		n := i % fixed_array_max_size
		sample := ba.Get(n)
		if sample.Num != uint32(n) {
			b.Errorf("got %v, want %v", sample.Num, i)
		}
	}
}

type Sample struct {
	//comment [200]byte
	Num uint32
	//status  bool
	mobileId [12]byte
	spId     [12]byte
	simId    [12]byte
	spId2    [12]byte
}
