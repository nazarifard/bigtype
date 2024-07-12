package sync

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

var (
	maxMapSize = func() int {
		size := os.Getenv("MAXSIZE")
		size = strings.Replace(size, "_", "", -1)
		if n, err := strconv.Atoi(size); err == nil {
			return n
		} else {
			return 100_000
		}
	}()
	threads    = 1 //runtime.NumCPU()
	iterations = 3 * maxMapSize

	cache     = make(map[int]int, 100)
	syncCache sync.Map
	mutex     sync.Mutex
	big       = NewMap[int, int](maxMapSize)
)

func memstat(name string, Operation func(int)) {
	runtime.GC()
	fmt.Printf("Engine: %s\t", name)
	var mem0, mem runtime.MemStats
	runtime.ReadMemStats(&mem0)
	Do(Operation)
	runtime.ReadMemStats(&mem)
	fmt.Printf("Sys: %4dMB, ", (mem.Sys-mem0.Sys)/1024/1024)
	fmt.Printf("Alloc: %4dMB, ", (mem.Alloc-mem0.Alloc)/1024/1024)
	fmt.Printf("HeapAlloc: %4dMB, ", (mem.HeapAlloc-mem0.HeapAlloc)/1024/1024)
	fmt.Printf("NumGC: %d\n", mem.NumGC-mem0.NumGC)
}

func Do(Operation func(int)) {
	wg := new(sync.WaitGroup)
	wg.Add(threads)
	for t := 0; t < threads; t++ {
		go func(id int) {
			for range iterations / threads {
				Operation(rand.Int() % maxMapSize)
			}
			wg.Done()
		}(t)
	}
	wg.Wait()
}

func Test_Memstats(t *testing.T) {
	memstat("BigMap   ", func(i int) { big.Update(i%maxMapSize, func(v int) int { return v + 1 }) })
	memstat("StdMap   ", func(i int) { mutex.Lock(); cache[i%maxMapSize] += 1; mutex.Unlock() })
	memstat("sync.Map ", func(i int) { elem, _ := syncCache.LoadOrStore(i, new(int32)); atomic.AddInt32(elem.(*int32), 1) })
}
