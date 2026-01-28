# bigtype
Bigtype provides big type simply for Golang. In addition to the simple types, Go offers three built-in complex types that are exceedingly useful: arrays, slices and maps.
Go standard built-in types work fine and can be enough often but not always. Maximum length of Go-lang arrays is 2GB. Then we can not define big array more than 2GB. on the other side Go maps are exteremly fast but uses a lot of pointers and when we have to define a big map, Golang standard maps are not good choice because use a lot of memory and finally application may be crashed. 
Bigtype provides 3 alternative types instead of built-in Go standard types.

# Example 
```go
import "github.com/nazarifard/bigtype"
import "fmt" 

func demo() error {
  bigArr:=bigtype.NewArray[string](200_000_000) 
  for k := range 200_000_000 {
	  a.set(k, fmt.Sprint(-k))
  }
  
  m:=bigtype.NewMap[int, string]()
  for k := range 200_000_000 {
	  m.set(k, fmt.Sprint(k))
  }
  
  for k, v := range m.Seq {
    if a[k] != "-" + v {
      return fmt.Errorf("bigType failed");
   }
  }
  return nil
}
```
# bigtype/sync
 bigtype/sync package provides thread-safe bigtypes. sync bigtypes implemeted as a scalable concurrent high performance big Go data type without engaging GC. 

# Benchmark
Golang sync.map and standard map respectively use almost 8x and 20x more memory. 
```bash
$ MAXSIZE=10_000_000 go test -bench=. -run=^$ -benchtime=30000000x
 goos: linux
 goarch: amd64
 pkg: github.com/nazarifard/bigtype
 cpu: Intel(R) Core(TM) i7-3537U CPU @ 2.00GHz
 Benchmark_MapStringString_Set-4                 30000000    1220.00 ns/op      6  B/op     0 allocs/op
 Benchmark_stdMap_StringString_Set-4             30000000     865.5 ns/op      47  B/op     1 allocs/op
 Benchmark_stdSyncMapStringString_Set-4          30000000    1843.00 ns/op     100 B/op     4 allocs/op
 
```
# Memstat
The report of memstat API shows the diffrence between bigtype and other similar solutions clearly.
Bigtype never uses GC and or HeapAlloc
```sh
$ MAXSIZE=10_000_000 go test -run=Memstats
 Engine: BigMap          Sys:    0MB, Alloc:    0MB, HeapAlloc:    0MB, NumGC: 0
 Engine: StdMap          Sys:  649MB, Alloc:  580MB, HeapAlloc:  580MB, NumGC: 1
 Engine: sync.Map        Sys: 1372MB, Alloc: 1453MB, HeapAlloc: 1453MB, NumGC: 2
```
# dgraph-io/cachebench
BigMap is an alias for bigtype.Map works as fast as fastest cache in the world, while uses zero heap memory allocations and never is engaging with GC! the result of benchmarks shows how Bigtype.Map (and also other underlayers bigtype.Array) work stable and reliable with 0 extra memory allocation.

```sh
BenchmarkCaches/BigMapZipfWrite-4          2416587       448.6 ns/op     0 B/op    0 allocs/op
BenchmarkCaches/SyncMapZipfWrite-4         1000000      1106.0 ns/op    63 B/op    4 allocs/op
BenchmarkCaches/RistrettoZipfWrite-4       1000000      2270.0 ns/op   128 B/op    3 allocs/op
BenchmarkCaches/GroupCacheZipfWrite-4      2632586       488.2 ns/op    47 B/op    2 allocs/op
                                                                                
BenchmarkCaches/BigMapOneKeyWrite-4        2918360       403.8 ns/op     0 B/op    0 allocs/op
BenchmarkCaches/SyncMapOneKeyWrite-4       2852132       422.9 ns/op    61 B/op    4 allocs/op
BenchmarkCaches/RistrettoOneKeyWrite-4     2994553       370.1 ns/op   128 B/op    3 allocs/op
BenchmarkCaches/GroupCacheOneKeyWrite-4    3638934       347.1 ns/op    45 B/op    3 allocs/op

BenchmarkCaches/BigMapOneKeyRead-4        10695729       113.9 ns/op     0 B/op    0 allocs/op
BenchmarkCaches/SyncMapOneKeyRead-4       33431958        31.6 ns/op     0 B/op    0 allocs/op
BenchmarkCaches/RistrettoOneKeyRead-4     12718990       345.0 ns/op    24 B/op    1 allocs/op
BenchmarkCaches/GroupCacheOneKeyRead-4     7333424       173.1 ns/op     0 B/op    0 allocs/op

BenchmarkCaches/BigMapZipfRead-4           3178783       369.3 ns/op     0 B/op    0 allocs/op
BenchmarkCaches/SyncMapZipfRead-4          6624446       154.2 ns/op     0 B/op    0 allocs/op
BenchmarkCaches/RistrettoZipfRead-4       12708687        94.2 ns/op    24 B/op    1 allocs/op
BenchmarkCaches/GroupCacheZipfRead-4       6086665       197.3 ns/op     0 B/op    0 allocs/op
```

# Internal modules 
Bigtype uses [fastape](github.com/nazarifard/fastape) and [marshaltap](github.com/nazarifard/marshaltap) as fastest marshaling solution and [syncpool](github.com/nazarifard/syncpool) for managing temporary ojects. as well as is using T1HA0 algorithm for hashing and encoding binary data.


# Applications
bigtype can be used widely in various application that need to load, handle and process a lot of data. including in-memory DBs, Cache solutions, Data migration tools, et cetera.
