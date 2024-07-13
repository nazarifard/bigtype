package bucket

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/nazarifard/bigtype/internal/addr"
)

// func Benchmark_Rand1(b *testing.B) {
// 	b.ResetTimer()
// 	for range b.N {
// 		offset := rand.Int63n(0xFFFFFFFFFF-1) % 63000 //0xFFFF
// 		Len := 1 + rand.Int31n(4)
// 		//copy(bucket.cellar[62000:], s[:4])
// 		_ = offset
// 		_ = Len
// 	}
// }

// func Benchmark_Copy1(b *testing.B) {
// 	bucket := NewBucket1()
// 	s := []byte("1234567890")
// 	b.ResetTimer()
// 	for range b.N {
// 		//offset := rand.Int63n(0xFFFFFFFFFF-1) % 63000 //0xFFFF
// 		//Len := 1 + rand.Int31n(4)
// 		copy(bucket.cellar[62000:], s[:4])
// 	}
// }

// func Benchmark_Insert1(b *testing.B) {
// 	bucket := NewBucket1()
// 	s := []byte("1234567890")
// 	b.ResetTimer()
// 	for range b.N {
// 		//index := rand.Int63n(0xFFFFFFFFFF - 1) //% 5200 //0xFFFF
// 		//Len := 1 + rand.Int31n(4)
// 		bucket.Write(63000, s[:4])
// 	}
// }

// func Benchmark_Defrag1(b *testing.B) {
// 	bucket := NewBucket1(0)
// 	//fmt.Println("_______________________________")
// 	s := []byte("1234567890")
// 	//m := make(map[int64]uint16)
// 	var offset uint16
// 	var ok bool
// 	del, defrag := 0, 0
// 	b.ResetTimer()
// 	for i := range b.N {
// 		index := rand.Int63n(0xFFFFFFFFFF-1) % 5200 //0xFFFF
// 		Len := 1 + rand.Int31n(4)
// 		if offset, ok = m[index]; ok && offset != NILOFFSET {
// 			bucket.Delete(offset)
// 			del++
// 			m[index] = 0
// 		}
// 		offset = bucket.Write(int(index), s[:Len])
// 		if offset == NILOFFSET {
// 			if bucket.deadSpace > 10+10+uint16(Len) {
// 				bucket.Defrag(0, m)
// 				defrag++
// 				offset = bucket.Write(int(index), s[:Len])
// 			} else {
// 				fmt.Println("............Bucket is really full")
// 				break
// 			}
// 		}
// 		if offset != NILOFFSET {
// 			m[index] = offset
// 		} else {
// 			fmt.Println(" !!!!!!!!!! NilOffset i", i, "free", bucket.FreeCell.Head.Len, "dead", bucket.deadSpace, "del", del)
// 		}
// 	}
// 	// fmt.Println("Ok free", bucket.FreeCell.Head.Len, "dead", bucket.DeadSpace, "del", del, "defrag", defrag)
// }

// func Benchmark_FindMax_Defrag1(b *testing.B) {
// 	bucket := NewBucket1()
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

func CheckAddr(b *Bucket1, at addr.AddressTable) {
	for index := range at.Len {
		offset := at.Get(index).Offset()
		if offset != NILOFFSET && string(b.Get(offset)) != fmt.Sprint(index) {
			panic("check read")
		}
	}
}

func Test_Bucket1(t *testing.T) {
	addrTable := addr.AddressTable{} //make(map[int64]uint16)
	addrTable.Expand(6000)
	addrTable.FixedSize = true

	b := NewBucket1(0)
	var index int
	for index = range 5000 {
		space, offset := b.Request(index, len([]byte(fmt.Sprint(index))))
		addrTable.Set(index, addr.NewAddrItem(0, offset))
		copy(space, []byte(fmt.Sprint(index))) //write
	}

	for index = range 5000 {
		ad := addrTable.Get(index)
		data := b.Get(ad.Offset())
		if string(data) != fmt.Sprint(index) {
			panic("set and get don't match")
		}
	}

	if err := b.CheckUp(addrTable); err != nil {
		panic(err)
	}
	for range 500 {
		r := int(1 + rand.Int63n(998))
		if addrTable.Get(r).Offset() != NILOFFSET {
			b.Delete(addrTable.Get(r).Offset())
			addrTable.Set(r, addr.NewAddrItem(0, NILOFFSET))
		}
	}
	if err := b.CheckUp(addrTable); err != nil {
		panic(err)
	}

	// m2 := make(map[int64]uint16)
	// for index, offset := range m {
	// 	if offset != 0 {
	// 		m2[index] = offset
	// 	}
	// }
	b.Defrag(&addrTable)
	if err := b.CheckUp(addrTable); err != nil {
		panic(err)
	}
	for index := 850; index < 900; index++ {
		if addrTable.Get(index).Offset() != NILOFFSET {
			b.Delete(addrTable.Get(index).Offset())
			addrTable.Set(index, addr.NewAddrItem(0, NILOFFSET))
		}
	}
	for index = 100; index < 600; index++ {
		if addrTable.Get(index).Offset() != NILOFFSET {
			b.Delete(addrTable.Get(index).Offset())
		}
		space, offset := b.Request(index, len([]byte(fmt.Sprint(index))))
		addrTable.Set(index, addr.NewAddrItem(0, offset))
		copy(space, []byte(fmt.Sprint(index))) //write
	}
	if err := b.CheckUp(addrTable); err != nil {
		panic(err)
	}
}

func Test_Insert1(t *testing.T) {
	bucket := NewBucket1(0)
	//bucket.CheckLen()
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
		space, offset := bucket.Request(int(index), len(s[:Len]))
		//	bucket.Checkup()
		if offset == NILOFFSET {
			bucket = NewBucket1(0)
			//		bucket.Checkup()
			copy(space, s[:Len])
			//offset = bucket.Write(int(index), s[:Len])
			//		bucket.Checkup()
		}
		if offset != NILOFFSET {
			m[index] = offset
		} else {
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		}
	}
}
