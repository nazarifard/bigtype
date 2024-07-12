# bigtype
Bigtype provides big type simply for Golang. In addition to the simple types, Go offers three built-in complex types that are exceedingly useful: arrays, slices and maps.
Go standard built-in types work fine and can be enough often but not always. Maximum length of Go-lang arrays is 2GB. Then we can not define big array more than 2GB. on the other side Go maps are exteremly fast but uses a lot of pointers and when we have to define a big map, Golang standard maps are not good choice because use a lot of memory and finally application may be crashed. 
Bigtype provides 3 alternative types instead of built-in Go standard types.

# Example 
```go
  import "github.com/nazarifard/bigtype"
   func demo() {
     //var a bigtype.Array[int]
     a:=bigtype.NewArray[string](500_000_000) 
     a.Set(123456789, "-123456789")
     fmt.Println(a.Get(123456789))

     //var m bigtype.Map[K,V]
     m:=bigtype.NewMap[int, string]()
     m.set(123, "123")
     m.set(456, "456")
     m.Range(func(key, value int) bool {
            print("{", key, ":", value, "}, ");
            return true
          })
    }
```
# bigtype/sync
 bigtype/sync package provides thread-safe bigtypes. sync bigtypes implemeted as a scalable concurrent high performance big Go data type without engaging GC. 

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
# Memstat
The report of memstat API shows the diffrence between bigtype and other similar solutions clearly.
Bigtype never uses GC and or HeapAlloc
```sh
$ MAXSIZE=10_000_000 go test -run=Memstats
 Engine: BigMap          Sys:    0MB, Alloc:    0MB, HeapAlloc:    0MB, NumGC: 0
 Engine: StdMap          Sys:  649MB, Alloc:  580MB, HeapAlloc:  580MB, NumGC: 1
 Engine: sync.Map        Sys: 1372MB, Alloc: 1453MB, HeapAlloc: 1453MB, NumGC: 2
```

# Internal modules 
Bigtype uses [fastape](github.com/nazarifard/fastape) and [marshaltap](github.com/nazarifard/marshaltap) as fastest marshaling solution and [syncpool](github.com/nazarifard/syncpool) for managing temporary ojects. as well as is using T1HA0 algorithm for hashing and encoding binary data.


# Applications
bigtype can be used widely in various application that need to load, handle and process a lot of data. including in-memory DBs, Cache solutions, Data migration tools, et cetera.





  
