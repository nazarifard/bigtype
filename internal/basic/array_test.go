package basic

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

//var _ = func() int { config.Load(); return 0 }()

const g_array_size = 12_000

func BenchmarkArray_Uint64(b *testing.B) {
	ba := NewArray[uint64](g_array_size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ba.Set(i%g_array_size, uint64(i))
	}
}

func Benchmark_Set_Array_Uint64(b *testing.B) {
	ba := NewArray[uint64]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ba.Set(i%g_array_size, uint64(i))
	}
}

func Benchmark_Set_MArray_Uint64(b *testing.B) {
	var ba_mArray_uint64 = NewArray[uint64](MTapeUint64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ba_mArray_uint64.Set(i%g_array_size, uint64(i))
	}
}

func Benchmark_Set_BigArray_String(b *testing.B) {
	ba := NewArray[string]()
	for i := 0; i < b.N; i++ {
		ba.Set(i%g_array_size, "abcdefghjklmnbop") //fmt.Sprintf("%0.8d", i))
	}
}
func Benchmark_Array_Set_MBig_Slice_String(b *testing.B) {
	ba := NewArray[string](MTapeString)
	for i := 0; i < b.N; i++ {
		ba.Set(i%g_array_size, "abcdefghjklmnbop") //fmt.Sprintf("%0.8d", i))
	}
}

//	func Test_MTree_Uint64(t *testing.T) {
//		//t := newBnfTree[int, ](g_array_size, HashT1ha0[string])
//		//t := newBnfTree[uint64](g_array_size, NoHash[uint64])
//		tree := newMTree[uint64, uint64, marshal.MUint64](g_array_size)
//		//seed:= rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
//		//seed:=rand.Uint64()
//		//m.InOrder(m.Root)
//		//	start := time.Now()
//		for i := 0; i < g_array_size/Batch; i++ {
//			for j := 0; j < Batch; j++ {
//				r := uint64(i*Batch + j)
//				//r := fmt.Sprint(i*Batch + j) //randStr(15)
//				//r := uint64(rand.Intn(Batch * 900))
//				tree.Insert(r, r)
//				//m.Insert(i*Batch + j)
//			}
//			//		fmt.Printf("i:%d, time:%v \n", i, time.Since(start))
//		}
//		//tree.Print()
//		fmt.Printf("\n{ ")
//		it := tree.Iterator()
//		for it.Next() {
//			fmt.Printf("{%d, %d}, ", it.Key(), it.Value())
//		}
//		fmt.Printf(" }\n")
//		//	start = time.Now()
//		for i := 0; i < g_array_size; i++ {
//			got, ok := tree.Search(uint64(i)) //fmt.Sprint(2 * 1000 * 1000))
//			if !ok {
//				t.Errorf("want:%d, got:%d", i, got)
//			}
//		}
//	}

func BenchmarkCacheGet(b *testing.B) {

	x := NewArray[int]()
	_ = x

	const items = 1 << 16
	ba := NewArray[[]byte](items)
	v := []byte("xyza")
	for i := 0; i < items; i++ {
		ba.Set(i, v)
	}
	fmt.Println("_____________________________________________________")
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		//k := 0 //[]byte("\x00\x00\x00\x00")
		for pb.Next() {
			for i := 0; i < items; i++ {
				_ = ba.Get(i)
			}
		}
	})
}

func BenchmarkCacheSet(b *testing.B) {
	const items = 1 << 16
	ba := NewArray[[]byte](items)
	fmt.Println("_____________________________________________________")
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		//k := 0 //[]byte("\x00\x00\x00\x00")
		v := []byte("xyza")
		for pb.Next() {
			for i := 0; i < items; i++ {
				ba.Set(i, v)
			}
		}
	})
}

func TestPerformance_BigArray_String(t *testing.T) {
	const size = g_array_size/2 - 1
	array := NewArray[[]byte](size)
	// for k := range size {
	// 	i := rand.Int31n(size)
	// 	s := fmt.Sprint(i)
	// 	array.Set(int(i), s) //fmt.Sprintf("%0.8d", i))
	// 	//	tree.Set(i, fmt.Sprint(-i))
	// 	_ = k
	// }

	start := time.Now()
	//sum := 0
	str := []byte("9876") //543210_9876543210_9876543210_9876543210_9876543210"
	ch := make(chan int, 4)
	f := func() {
		sum := 0
		for i := range 10 * size {
			//i := rand.Int31n(size)
			//s := fmt.Sprint(i)
			array.Set(i%size, str) //int(i%size), str[:rand.Int31n(32)])
			sum += len(str)
		}
		ch <- sum
	}
	go f()
	// go f()
	// go f()
	// go f()
	sum := <-ch //+ <-ch + <-ch + <-ch

	duration := time.Since(start).Seconds()
	fmt.Printf("\nperformance:%.2f MB", float64(sum)/(duration*1024*1024))

	// start = time.Now()
	// //sum := 0
	// str = "987654"
	// ch = make(chan int, 4)
	// f = func() {
	// 	sum := 0
	// 	for i := range size {
	// 		//i := rand.Int31n(size)
	// 		//s := fmt.Sprint(i)
	// 		str = array.Get(int(i))
	// 		sum += len(str)
	// 	}
	// 	ch <- sum
	// }
	// go f()
	// go f()
	// go f()
	// go f()
	// sum = <-ch + <-ch + <-ch + <-ch

	// duration = time.Since(start).Seconds()
	// fmt.Printf("\nGet performance:%.2f MB", float64(sum)/(duration*1024*1024))

}

func Test_BigArray_String(t *testing.T) {
	arr := NewArray[string](g_array_size)
	for k := range g_array_size {
		i := rand.Int31n(g_array_size)
		arr.Set(int(i), fmt.Sprint(i)) //fmt.Sprintf("%0.8d", i))
		//	tree.Set(i, fmt.Sprint(-i))
		_ = k
	}
	for i := 0; i < g_array_size; i++ {
		i := rand.Int31n(g_array_size)
		arr.Set(int(i), fmt.Sprint(-i))
	}
	for i := 1; i < g_array_size; i++ {
		j := arr.Get(i) //fmt.Sprintf("%0.8d", i))
		if j != fmt.Sprint(-i) {
			//fmt.Println(i, j)
			_ = j
		}
	}
}

func Test_BigArray_Uint64(t *testing.T) {
	arr := NewArray[int]()
	for i := 0; i < g_array_size; i++ {
		arr.Set(i, i) //fmt.Sprintf("%0.8d", i))
	}
	for i := 0; i < g_array_size; i++ {
		arr.Set(i, -i) //fmt.Sprintf("%0.8d", i))
	}
	for i := 0; i < g_array_size; i++ {
		j := arr.Get(i) //fmt.Sprintf("%0.8d", i))
		if j != -i {
			fmt.Println(i, j)
		}
	}
}
