package basic

//TODO needs more tests yet
// import (
// 	"math/rand"
// 	"testing"
// )
//
// func Test_FsMap_Uint640(t *testing.T) {
// 	fsMap := NewHard()
// 	//b.ResetTimer()
// 	for i := 0; i < 3000; i++ {
// 		k := rand.Uint64() //uint64(i % g_map_size)
// 		j := rand.Uint64() //uint64(i % g_map_size)
// 		fsMap.Set(uint64(k), int(j))
// 		//fmt.Println("aaaaaaaaaaaaaaaa")
// 	}
// 	// for i := 0; i < fsMap.Len(); i++ {
// 	// 	litter.Dump(fsMap.Array.Get(i))
// 	// }
// }
//
// func Benchmark_FsMap_Uint64(b *testing.B) {
// 	fsMap := NewDisk()
// 	for i := 0; i < b.N; i++ {
// 		k := rand.Uint64()
// 		j := rand.Uint64()
// 		fsMap.Set(uint64(k), int(j))
// 	}
// }
// func Benchmark_FsMap_Uint64_Get(b *testing.B) {
// 	fsMap := NewHard()
// 	//b.ResetTimer()
// 	for i := 0; i < 300_000; i++ {
// 		k := rand.Uint64() //uint64(i % g_map_size)
// 		j := rand.Uint64() //uint64(i % g_map_size)
// 		fsMap.Set(uint64(k), int(j))
// 		//fmt.Println("aaaaaaaaaaaaaaaa")
// 	}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		k := rand.Uint64() //uint64(i % g_map_size)
// 		_, _ = fsMap.Get(uint64(k))
// 	}
// 	// for i := 0; i < fsMap.Len(); i++ {
// 	// 	litter.Dump(fsMap.Array.Get(i))
// 	// }
// }
//
// func Test_FsMap_Uint64(t *testing.T) {
// 	fsMap := NewHard()
// 	fsMap.Set(0x7766554433221100, 111)
// 	fsMap.Set(0x8866554433221100, 222)
// 	fsMap.Set(0x7766554433221100, 333)
// 	fsMap.Set(0x7766554433221100, 0x76543210)
// 	fsMap.Set(0x7766554433221100, 0x76543210)
// 	fsMap.Set(0x7766554433221100, 0x76543210)
// 	fsMap.Set(0x7766554433221100, 0x76543210)
// 	//fmt.Println(fsMap.Len())
// }
