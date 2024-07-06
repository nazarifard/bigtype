# bigtype
scalable concurrent high performance big Go data type without engaging with GC 

bigtype provides big type simply for Golang. Go-lang supports arrays/slices and maps by defaults. They are work fine and can be enough often but not always. Maximum length of Go-lang arrays is 2GB. Then we can not define big array more than 2GB. on the other side Go maps are exteremly fast but uses a lot of pointers and when we have to define a big map, Golang standard maps are not good choice because use a lot of memory and finally application may be crashed. Using an appropriate Cache is typical solution to solve the problem. But as I checked all availiable cache solution finally using a type of standard go map. Then they can not resolve the problem but only can reduce the impact or improve the result a little.  
bigtype main idea is fundamentaly diffrent. Every thing is implemented on array without using of Go standard maps. Therefore bigtype never engaging with GC while uses minimum required memory. 

# Example 
```go
  import "github.com/nazarifard/bigtype"
   func demo() {
     //var a bigtype.Array[int]
     //var m bigtype.Map[K,V]
     //var b bigtype.BitArray
     //var c bigtype.Cache
     a:=bigtype.NewArray[string](500_000_000) 
     a.Set(123456789, "-123456789")
     fmt.Println(a.Get(123456789))
  }
```
# Benchmark
Golang sync.map and standard map respectively use almost 5x and 10x more memory. 
```bash
$ MAXSIZE=10_000_000 go test -bench=. -run=^$ -benchtime=30000000x
 goos: linux
 goarch: amd64
 pkg: github.com/nazarifard/bigtype
 cpu: Intel(R) Core(TM) i7-3537U CPU @ 2.00GHz
 Benchmark_MapStringString_Set-4                 30000000    1220.00 ns/op     10  B/op     0 allocs/op
 Benchmark_stdMap_StringString_Set-4             30000000     865.5 ns/op      47  B/op     1 allocs/op
 Benchmark_stdSyncMapStringString_Set-4          30000000    1843.00 ns/op     100 B/op     4 allocs/op
```

Other solutions like freecash, bigcache, fastcache and similar tools don't have better situation. For example Ristretto as a best high-performance memory-bound cache based on [graph-io/cachebenchmark](https://github.com/dgraph-io/benchmarks/tree/master/cachebench) benchmark result also uses a lot of memory. The reason is simple, they are finally using Go standard map at backend then can not solve the memory usage problem fundamentally. However Go standard map is exteremly fast and all solutions based on Go map will be fast also.
```sh
 $ go test -bench=. -cpu=4 -timeout 0
 goos: windows
 goarch: amd64
 pkg: github.com/nazarifard/bigtype/cachebench
 cpu: Intel(R) Core(TM) i7-3537U CPU @ 2.00GHz
 BenchmarkCaches/RistrettoZipfWrite-4      3905212     276.5 ns/op      128 B/op          3 allocs/op
 BenchmarkCaches/BigtypeZipfWrite-4        2784415     396.0 ns/op        0 B/op          0 allocs/op

 BenchmarkCaches/RistrettoOneKeyWrite-4    3368662     373.8 ns/op      128 B/op          3 allocs/op
 BenchmarkCaches/BigtypeOneKeyWrite-4      1915390     636.4 ns/op        0 B/op          0 allocs/op
```
# Memstat
Golang memstat API report shows the diffrence between bigtype and other similar solutions clearly.
Bigtype never uses GC and or HeapAlloc
```sh
$ MAXSIZE=10_000_000 go test -run=MemStat
  Engine: BigtypeMap    Sys:    0MB, Alloc:    0MB, HeapAlloc:    0MB, NumGC: 0
  Engine: StdMap        Sys:  649MB, Alloc:  580MB, HeapAlloc:  580MB, NumGC: 1
  Engine: sync.Map      Sys:  888MB, Alloc: 1380MB, HeapAlloc: 1380MB, NumGC: 1
```

# Internal modules 
Bigtype uses [fastape](github.com/nazarifard/fastape) and [marshaltap](github.com/nazarifard/marshaltap) as fastest marshaling solution and [syncpool](github.com/nazarifard/syncpool) for managing temporary ojects. as well as is using T1HA0 algorithm for hashing and encoding binary data.


# Applications
bigtype can be used widely in various application that need to load, handle and process a lot of data. including in-memory DBs, Cache solutions, Data migration tools, et cetera.





  
