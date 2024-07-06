package basic

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func Benchmark_Map_StringString(b *testing.B) {
	bigmap := NewMap[string, string]() //(0, MTapeString, MTapeUint64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := fmt.Sprint(i % 2_000_000)
		bigmap.Set(s, s) //fmt.Sprintf("%0.8d", k))
	}
	// fmt.Println("bigmap.Len(): ", bigmap.Len())
}

func Benchmark_Map_String(b *testing.B) {
	bigmap := NewMap[string, uint64](0, MTapeString, MTapeUint64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bigmap.Set(fmt.Sprint(i%2_000_000), uint64(i)) //fmt.Sprintf("%0.8d", k))
	}
	// fmt.Println("bigmap.Len(): ", bigmap.Len())
}

func TestBigMap(t *testing.T) {
	const SIZE = 10_000
	bigmap := NewMap[string, string]() //(0, MTapeString, MTapeUint64)

	start := time.Now()
	for range uint64(SIZE) {
		i := rand.Uint64() % SIZE
		bigmap.Set(fmt.Sprint(i), fmt.Sprint(i))
	}
	fmt.Printf("\n insert time: %v", time.Since(start))

	start = time.Now()
	for range uint64(SIZE) {
		i := rand.Uint64() % SIZE
		bigmap.Set(fmt.Sprint(i), fmt.Sprint(i)+"1")
	}
	fmt.Printf("\n update time: %v", time.Since(start))

	start = time.Now()
	for range uint64(SIZE) {
		i := rand.Uint64() % SIZE
		k := fmt.Sprint(i)
		v, ok := bigmap.Get(k)
		// _, _ = j, ok
		// fmt.Println(i, j, ok)
		if ok && k != v && k != v+"1" {
			t.Errorf("i!=-j %v!=%v", k, v)
			os.Exit(1)
		}
	}
	fmt.Printf("\n search time: %v", time.Since(start))

	fmt.Println()
}

func TestBigMapUint64(t *testing.T) {
	const SIZE = 10_000
	bigmap := NewMap[uint64, uint64](0, MTapeUint64, MTapeUint64)

	start := time.Now()
	for range uint64(SIZE) {
		i := rand.Uint64() % SIZE
		bigmap.Set(i, 10*i)
	}
	fmt.Printf("\n insert time: %v", time.Since(start))

	start = time.Now()
	for i := range uint64(SIZE) {
		bigmap.Set(i, i+1)
	}
	fmt.Printf("\n update time: %v", time.Since(start))

	start = time.Now()
	for range uint64(SIZE) {
		i := rand.Uint64() % SIZE
		j, ok := bigmap.Get(i)
		//fmt.Println(i, j, ok)
		if !ok || j != i+1 {
			t.Errorf("i!=-j %v!=%v", i, j)
			os.Exit(1)
		}
	}
	fmt.Printf("\n search time: %v", time.Since(start))

	fmt.Println()
}

func TestMapPrint(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("001", 1)
	m.Set("002", 2)
	m.Set("003", 3)
	m.Set("002", -2)
	//m.print()
	m.Range(func(key string, value int) bool {
		fmt.Println("key:", key, "value:", value, ", ")
		return true
	})

	//}
}
