package bucket

import (
	"fmt"
	"math/rand"
	"testing"
)

func Benchmark_Rand(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		offset := rand.Int63n(0xFFFFFFFFFF-1) % 63000 //0xFFFF
		Len := 1 + rand.Int31n(4)
		//copy(bucket.cellar[62000:], s[:4])
		_ = offset
		_ = Len
	}
}

func Benchmark_Copy(b *testing.B) {
	bucket := NewBucket2()
	s := []byte("1234567890")
	b.ResetTimer()
	for range b.N {
		//offset := rand.Int63n(0xFFFFFFFFFF-1) % 63000 //0xFFFF
		//Len := 1 + rand.Int31n(4)
		copy(bucket.cellar[62000:], s[:4])
	}
}

func Benchmark_Insert(b *testing.B) {
	bucket := NewBucket2()
	s := []byte("1234567890")
	b.ResetTimer()
	for range b.N {
		//index := rand.Int63n(0xFFFFFFFFFF - 1) //% 5200 //0xFFFF
		//Len := 1 + rand.Int31n(4)
		bucket.Write(63000, s[:4])
	}
}

func Benchmark_Defrag(b *testing.B) {
	bucket := NewBucket2()
	//fmt.Println("_______________________________")
	s := []byte("1234567890")
	m := make(map[int64]uint16)
	var offset uint16
	var ok bool
	del, defrag := 0, 0
	b.ResetTimer()
	for i := range b.N {
		index := rand.Int63n(0xFFFFFFFFFF-1) % 5200 //0xFFFF
		Len := 1 + rand.Int31n(4)
		if offset, ok = m[index]; ok && offset != NILOFFSET {
			bucket.Delete(offset)
			del++
			m[index] = 0
		}
		offset = bucket.Write(int(index), s[:Len])
		if offset == NILOFFSET {
			if bucket.DeadSpace > 10+10+uint16(Len) {
				bucket.Defrag(m)
				defrag++
				offset = bucket.Write(int(index), s[:Len])
			} else {
				fmt.Println("............Bucket is really full")
				break
			}
		}
		if offset != NILOFFSET {
			m[index] = offset
		} else {
			fmt.Println(" !!!!!!!!!! NilOffset i", i, "free", bucket.FreeCell.Head.Len, "dead", bucket.DeadSpace, "del", del)
		}
	}
	// fmt.Println("Ok free", bucket.FreeCell.Head.Len, "dead", bucket.DeadSpace, "del", del, "defrag", defrag)
}

// func Benchmark_FindMax_Defrag(b *testing.B) {
// 	bucket := NewBucket2()
// 	//fmt.Println("_______________________________")
// 	s := []byte("1234567890")
// 	m := make(map[int64]uint16)
// 	var offset uint16
// 	var ok bool
// 	del := 0
// 	b.ResetTimer()
// 	lock := false
// 	var Len int32
// 	var index int64
// 	for i := range b.N {
// 		//fmt.Println(i)
// 		if !lock {
// 			index = rand.Int63n(0xFFFFFFFFFF-1) % 5200 //0xFFFF
// 			Len = 1 + rand.Int31n(4)
// 			if offset, ok = m[index]; ok && offset != NILOFFSET {
// 				bucket.Delete(offset)
// 				del++
// 				m[index] = 0
// 			}
// 			offset = bucket.Write(int(index), s[:Len])
// 		}
// 		if offset == NILOFFSET && bucket.DeadSpace > uint16(Len)+10 {
// 			c, ok := bucket.FindMaxDead()
// 			cLen := c.Head.Len
// 			if ok && c.Head.Len > 10+10+uint16(Len) &&
// 				c.Head.Len > bucket.FreeCell.Head.Len {
// 				bucket.SetFreeCell(c)
// 			} else {
// 				bucket.Defrag(m)
// 			}
// 			offset = bucket.Write(int(index), s[:Len])
// 			if offset == NILOFFSET {
// 				fmt.Println(cLen, c.Head.Len, "AAAAAAAAAAAAAAAA")
// 				lock = true
// 				continue
// 			}
// 		}
// 		// 	c, ok = bucket.FindMaxDead()
// 		// 	if ok &&
// 		// 		c.Head.Len > uint16(10+Len) &&
// 		// 		c.Head.Len > bucket.FreeCell.Head.Len {
// 		// 		bucket.SetFreeCell(c)
// 		// 		//continue
// 		// 	} else {
// 		// 		bucket.Defrag(m)
// 		// 		//bucket.FindMaxDead()
// 		// 		//if bucket.FreeCell.Len >
// 		// 		fmt.Println("  Defrag", i, "free", bucket.FreeCell.Head.Len, "dead", bucket.DeadSpace, "del", del)
// 		// 		//break
// 		// 	}
// 		// 	offset = bucket.Write(int(index), s[:Len])
// 		// 	//break
// 		// 	//offset = bucket.Write(int(index), s[:Len])
// 		// }
// 		if offset != NILOFFSET {
// 			m[index] = offset
// 		} else {
// 			fmt.Println(" NilOffset i", i, "free", bucket.FreeCell.Head.Len, "dead", bucket.DeadSpace, "del", del)
// 			//break
// 		}
// 	}
// }

func Test_Bucket(t *testing.T) {
	m := make(map[int64]uint16)
	b := NewBucket2()
	var i int
	for i = range 1000 {
		a := b.Write(i, []byte(fmt.Sprint(i)))
		m[int64(i)] = a
	}
	b.CheckLen()
	b.CheckAddr(m)

	for range 500 {
		r := int64(1 + rand.Int63n(998))
		if m[r] != 0 {
			b.Delete(m[r])
			m[r] = 0
		}
	}
	b.CheckLen()
	b.CheckAddr(m)

	m2 := make(map[int64]uint16)
	for index, offset := range m {
		if offset != 0 {
			m2[index] = offset
		}
	}
	b.Defrag(m)

	b.CheckLen()
	b.CheckAddr(m)

	for i = 850; i < 900; i++ {
		//b.Delete(m[int64(i)])
		r := int64(i)
		if m[r] != 0 {
			b.Delete(m[r])
			m[r] = 0
		}
	}
	for i = 100; i < 600; i++ {
		if _, ok := m[int64(i)]; ok && m[int64(i)] != NILOFFSET {
			b.Delete(m[int64(i)])
		}
		a := b.Write(i, []byte(fmt.Sprint(i)))
		m[int64(i)] = a
	}
	b.CheckLen()
	b.CheckAddr(m)
}

func Test_Insert(t *testing.T) {
	bucket := NewBucket2()
	bucket.CheckLen()
	s := []byte("1234567890")
	m := make(map[int64]uint16)
	var offset uint16
	var ok bool
	for range 1_000 {
		//	fmt.Println(i)
		index := rand.Int63n(0xFFFFFFFFFF-1) % 0xFFFF
		Len := 1 + rand.Int31n(4)
		if offset, ok = m[index]; ok && m[index] != 0 {
			//bucket.Checkup()
			bucket.Delete(offset)
			//bucket.Checkup()
			m[index] = 0
		}
		offset = bucket.Write(int(index), s[:Len])
		//	bucket.Checkup()
		if offset == NILOFFSET {
			bucket = NewBucket2()
			//		bucket.Checkup()
			offset = bucket.Write(int(index), s[:Len])
			//		bucket.Checkup()
		}
		if offset != NILOFFSET {
			m[index] = offset
		} else {
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		}
	}
}
