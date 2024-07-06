// Developer: bahador.nazarifard@gmail.com
package bigtype

import (
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"unsafe"
)

var maxSize = func() int {
	n, err := strconv.Atoi(strings.Replace(os.Getenv("MAXSIZE"), "_", "", -1))
	if err == nil {
		return n
	} else {
		return 300_000 //default
	}
}()

func Benchmark_BitArray_Set(b *testing.B) {
	repo := NewArray[bool](maxSize)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := rand.Int31n(int32(maxSize))
			repo.Set(int(i), i&1 == 1)
		}
	})
}
func Benchmark_BitArray_Get(b *testing.B) {
	repo := NewArray[bool](maxSize)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := rand.Int31n(int32(maxSize))
			bit := repo.Get(int(i))
			_ = bit
		}
	})
}

func Benchmark_ArrayString_Get(b *testing.B) {
	repo := arrayStringSetup()
	//fmt.Println("------------------------------------------")
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := rand.Int31n(int32(maxSize))
			v := repo.Get(int(i))
			_ = v
		}
	})
}

func Benchmark_ArrayString_Set(b *testing.B) {
	alphabet := []byte("1234567890abcdef")
	//repo := arrayStringSetup()
	repo := NewArray[[]byte](maxSize)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := rand.Int31n(int32(maxSize))
			le := 1 + rand.Int31n(3)
			repo.Set(int(i), alphabet[:le])
		}
	})
}
func Benchmark_MapIntInt_Set(b *testing.B) {
	repo := NewMap[int, int](maxSize)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := int(1 + rand.Int31n(int32(maxSize)))
			repo.Set(n, n)
		}
	})
}
func Benchmark_MapIntInt_Get(b *testing.B) {
	repo := mapIntIntSetup()
	//fmt.Println("------------------------------------------")
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := int(1 + rand.Int31n(int32(maxSize)))
			v, ok := repo.Get(n)
			_, _ = v, ok
		}
	})
}

func Benchmark_MapIntString_Set(b *testing.B) {
	alphabet := []byte("1234567890abcdef")
	repo := NewMap[uint32, []byte](maxSize)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := rand.Int31n(int32(maxSize))
			le := 1 + rand.Int31n(3)
			repo.Set(uint32(i), alphabet[:le])
		}
	})
}

func Benchmark_MapIntString_Get(b *testing.B) {
	repo := mapIntStringSetup()
	//fmt.Println("------------------------------------------")
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := rand.Int31n(int32(maxSize))
			v, ok := repo.Get(uint32(i))
			_, _ = v, ok
		}
	})
}

func arrayStringSetup() Array[[]byte] {
	repo := NewArray[[]byte](maxSize)
	alphabet := []byte("1234567890abcdef")
	cpu := runtime.NumCPU()
	wg := sync.WaitGroup{}
	wg.Add(cpu)
	for range cpu {
		go func() {
			defer wg.Done()
			for range maxSize / runtime.NumCPU() {
				i := rand.Int31n(int32(maxSize))
				//s := fmt.Sprint(i)
				le := 1 + rand.Int31n(3)
				repo.Set(int(i), alphabet[:le])
			}
		}()
	}
	wg.Wait()
	return repo
}

func mapIntStringSetup() Map[uint32, []byte] {
	repo := NewMap[uint32, []byte](maxSize)
	alphabet := []byte("1234567890abcdef")
	cpu := runtime.NumCPU()
	wg := sync.WaitGroup{}
	wg.Add(cpu)
	for range cpu {
		go func() {
			defer wg.Done()
			for range maxSize / runtime.NumCPU() {
				i := rand.Int31n(int32(maxSize))
				//s := fmt.Sprint(i)
				le := 1 + rand.Int31n(3)
				repo.Set(uint32(i), alphabet[:le])
			}
		}()
	}
	wg.Wait()
	//mapIntString = repo
	return repo
}

func mapStringStringSetup() Map[string, []byte] {
	repo := NewMap[string, []byte](maxSize)
	alphabet := []byte("1234567890abcdef")
	cpu := runtime.NumCPU()
	wg := sync.WaitGroup{}
	wg.Add(cpu)
	for id := range cpu {
		go func(c int) {
			defer wg.Done()
			for i := range maxSize / cpu {
				//i := id*cpu +i //rand.Int31n(int32(hintSize))
				//s := fmt.Sprint(i)
				le := 1 + rand.Int31n(3)

				n := c*cpu + i
				p := (*byte)(unsafe.Pointer(&n))
				nStr := unsafe.String(p, unsafe.Sizeof(n))
				repo.Set(nStr, alphabet[:le])
				_ = i
			}
		}(id)
	}
	wg.Wait()
	return repo
}

func mapIntIntSetup() Map[int, int] {
	repo := NewMap[int, int](maxSize)
	cpu := runtime.NumCPU()
	wg := sync.WaitGroup{}
	wg.Add(cpu)
	for id := range cpu {
		go func(c int) {
			defer wg.Done()
			for range maxSize / cpu {
				n := int(rand.Int31n(int32(maxSize)))
				repo.Set(n, -n)
			}
		}(id)
	}
	wg.Wait()
	return repo
}

func Benchmark_MapIntInt(b *testing.B) {
	repo := NewMap[int, int](maxSize)
	cpu := runtime.NumCPU()
	wg := sync.WaitGroup{}
	wg.Add(cpu)
	for id := range cpu {
		go func(c int) {
			defer wg.Done()
			for range b.N / cpu {
				n := int(rand.Int31n(int32(maxSize)))
				_, _ = repo.Get(n)
				repo.Set(n, -n)
			}
		}(id)
	}
	wg.Wait()
}
func Benchmark_stdMapIntInt(b *testing.B) {
	repo := make(map[int]int, maxSize)
	cpu := runtime.NumCPU()
	wg := sync.WaitGroup{}
	mut := &sync.RWMutex{}
	wg.Add(cpu)
	for id := range cpu {
		go func(c int) {
			defer wg.Done()
			for range b.N / cpu {
				n := int(rand.Int31n(int32(maxSize)))
				mut.Lock()
				_ = repo[n]
				repo[n] = -n
				mut.Unlock()
			}
		}(id)
	}
	wg.Wait()
}

func Benchmark_MapStringString_Set(b *testing.B) {
	alphabet := []byte("1234567890abcdef")
	repo := NewMap[string, []byte](maxSize)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var n [4]byte
		for pb.Next() {
			le := 1 + rand.Int31n(3)

			p := (*int32)(unsafe.Pointer(&n))
			*p = rand.Int31n(int32(maxSize))
			q := unsafe.SliceData(n[:])
			nStr := unsafe.String(q, 4)

			repo.Set(nStr, alphabet[:le])
		}
	})
}
func Benchmark_MapStringString_Get(b *testing.B) {
	repo := mapStringStringSetup()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var n [4]byte
		for pb.Next() {
			p := (*int32)(unsafe.Pointer(&n))
			*p = rand.Int31n(int32(maxSize))
			q := unsafe.SliceData(n[:])
			nStr := unsafe.String(q, 4)
			v, ok := repo.Get(nStr)
			_, _ = v, ok
		}
	})
}

func setupSyncMapStringString() *sync.Map {
	var repo sync.Map
	alphabet := []byte("1234567890abcdef")
	cpu := runtime.NumCPU()
	wg := sync.WaitGroup{}
	wg.Add(cpu)
	for id := range cpu {
		go func(c int) {
			defer wg.Done()
			for i := range maxSize / cpu {
				n := c*maxSize/cpu + i
				p := (*byte)(unsafe.Pointer(&n))
				nStr := unsafe.String(p, unsafe.Sizeof(n))

				le := 1 + rand.Int31n(3)
				repo.Store(nStr, alphabet[:le])
			}
		}(id)
	}
	wg.Wait()
	return &repo
}

func Benchmark_stdSyncMapStringString_Set(b *testing.B) {
	alphabet := []byte("1234567890abcdef")
	var repo sync.Map //:= setupSyncMapStringString()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			le := 1 + rand.Int31n(3)
			n := 1 + rand.Int31n(int32(maxSize))
			p := (*byte)(unsafe.Pointer(&n))
			nStr := unsafe.String(p, unsafe.Sizeof(n))
			repo.Store(nStr, alphabet[:le])
		}
	})
}
func Benchmark_stdSyncMapStringString_Get(b *testing.B) {
	repo := setupSyncMapStringString()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var data []byte
		for pb.Next() {
			n := 1 + rand.Int31n(int32(maxSize))
			p := (*byte)(unsafe.Pointer(&n))
			nStr := unsafe.String(p, unsafe.Sizeof(n))
			v, ok := repo.Load(nStr)
			if ok {
				value, ok := v.([]byte)
				if ok {
					data = value
				}
			}
		}
		_ = data
	})
}

// func Benchmark_stdMap_StringString_Set(b *testing.B) {
// 	alphabet := []byte("1234567890abcdef")
// 	repo := make(map[string]string)
// 	mutex := &sync.RWMutex{}
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			le := 1 + rand.Int31n(3)
// 			n := 1 + rand.Int31n(int32(maxSize))
// 			p := (*byte)(unsafe.Pointer(&n))
// 			nStr := unsafe.String(p, unsafe.Sizeof(n))
// 			mutex.Lock()
// 			repo[nStr] = unsafe.String(&alphabet[0], le)
// 			mutex.Unlock()
// 		}
// 		if repo[fmt.Sprint(rand.Int31n(int32(maxSize)))] == "1230000" {
// 			print(repo["123"])
// 		}
// 	})
// }

// func setupStdMapStringString() map[string]string {
// 	repo := make(map[string]string)
// 	alphabet := []byte("1234567890abcdef")
// 	cpu := runtime.NumCPU()
// 	wg := sync.WaitGroup{}
// 	wg.Add(cpu)
// 	mutext := &sync.RWMutex{}
// 	for id := range cpu {
// 		go func(c int) {
// 			defer wg.Done()
// 			for i := range maxSize / cpu {
// 				n := c*maxSize/cpu + i
// 				p := (*byte)(unsafe.Pointer(&n))
// 				nStr := unsafe.String(p, unsafe.Sizeof(n))
// 				le := 1 + rand.Int31n(3)
// 				mutext.Lock()
// 				repo[nStr] = unsafe.String(&alphabet[0], le)
// 				mutext.Unlock()
// 			}
// 		}(id)
// 	}
// 	wg.Wait()
// 	return repo
// }

// func Benchmark_stdMap_StringString_Get(b *testing.B) {
// 	repo := setupStdMapStringString()
// 	mutex := &sync.RWMutex{}
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			n := 1 + rand.Int31n(int32(maxSize))
// 			p := (*byte)(unsafe.Pointer(&n))
// 			nStr := unsafe.String(p, unsafe.Sizeof(n))
// 			mutex.RLock()
// 			v, ok := repo[nStr]
// 			mutex.RUnlock()
// 			_, _ = v, ok
// 		}
// 	})
// }

// func Benchmark_Hash(b *testing.B) {
// 	h := hash.NewHash[string]() //hash.HashBuilder[string]()
// 	var n [4]byte
// 	for range b.N {
// 		p := (*int32)(unsafe.Pointer(&n))
// 		*p = rand.Int31n(int32(maxSize))
// 		q := unsafe.SliceData(n[:])
// 		nStr := unsafe.String(q, 4)
// 		// n := 1 + rand.Int31n(int32(hintSize))
// 		// p := (*byte)(unsafe.Pointer(&n))
// 		// nStr := unsafe.String(p, unsafe.Sizeof(n))
// 		h.Hash(nStr)
// 	}
// }

// func s2b(s string) []byte {
// 	p := unsafe.StringData(s)
// 	b := unsafe.Slice(p, len(s))
// 	return b
// }

// func Benchmark_Hash2(b *testing.B) {
// 	//h := hash.HashBuilder[string]()
// 	for range b.N {
// 		n := rand.Int31n(int32(hintSize))
// 		p := (*byte)(unsafe.Pointer(&n))
// 		nStr := unsafe.Slice(p, unsafe.Sizeof(n))
// 		hash.T1ha0(nStr, 0x12345678)
// 	}
// 	//print(x)
// }
