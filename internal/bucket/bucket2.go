package bucket

// import (
// 	"fmt"
// 	"unsafe"
// )

// type BucketHeader2 struct {
// 	DeadSpace uint16
// 	FreeCell  Cell2
// 	Len       int
// }

// type Bucket2 struct { //type Bucket [BucketSize]byte
// 	BucketHeader2
// 	cellar [BucketSize]byte
// }

// //var tmpBucket = Bucket{}

// func (b *Bucket2) MakeCell(offset, Len uint16) (c Cell2, ok bool) {
// 	if int(offset)+int(Len) > len(b.cellar) {
// 		return c, false
// 	}
// 	c.headPtr = (*Head)(unsafe.Pointer(&b.cellar[offset])) //cellar[0] is forbiden
// 	offset += uint16(unsafe.Sizeof(*c.headPtr))
// 	c.bodyPtr = &b.cellar[offset]
// 	c.Tail = (*Tail)(unsafe.Pointer(&b.cellar[offset+Len]))
// 	return c, true
// }

// func (b *Bucket2) Reset() {
// 	b.BucketHeader2 = BucketHeader2{}
// 	b.FreeCell, _ = b.MakeCell(1, uint16(len(b.cellar)-1-10))
// 	b.FreeCell.headPtr.Status = Dead
// 	b.FreeCell.headPtr.SetIndex(FreeCellIndex)
// 	b.FreeCell.headPtr.Len = uint16(len(b.cellar) - 1)
// 	b.FreeCell.Tail.Tlen = b.FreeCell.headPtr.Len
// }

// func NewBucket2() *Bucket2 {
// 	bucket := new(Bucket2)
// 	bucket.Reset()
// 	//fmt.Printf("\ntaked %v bytes memory for a new data bucket\n", cap(bucket.cellar))
// 	// if log.VerboseMode {
// 	// 	log.Logger.Info(fmt.Sprintf("taked %v bytes memory for a new data bucket", unsafe.Sizeof(*bucket)))
// 	// }
// 	return bucket
// }

// //	func (bucket *Bucket) Bytes() []byte {
// //		return bucket.cellar[:]
// //	}
// //
// //	func (b *Bucket) Set(key int, value []byte) (offset uint16) {
// //		return b.write(key, value)
// //	}
// func (b *Bucket2) Write(index int, data []byte) (offset uint16) {
// 	freeOffset := b.Offset(b.FreeCell)
// 	freeLen := b.FreeCell.headPtr.Len
// 	b.FreeCell.bodyPtr = &b.cellar[freeOffset+8] // : freeOffset+8-10+freeLen]

// 	if 8+len(data)+2 > int(freeLen)-8-2 {
// 		return NILOFFSET
// 	}
// 	if index > MaxValidIndex {
// 		panic("index is too large more than 5 bytes.")
// 	}

// 	newCell, _ := b.MakeCell(freeOffset, uint16(len(data)))
// 	newCell.headPtr.SetIndex(index)
// 	newCell.headPtr.Status = Live
// 	bs := unsafe.Slice(newCell.bodyPtr, len(data))
// 	copy(bs, data)
// 	newCell.headPtr.Len = uint16(8 + len(bs) + 2)
// 	newCell.Tail.Tlen = newCell.headPtr.Len
// 	offset = freeOffset //b.Offset(newCell)

// 	nextOffset := freeOffset + newCell.headPtr.Len
// 	var ok bool
// 	b.FreeCell, ok = b.MakeCell(nextOffset, freeLen-10-newCell.headPtr.Len)
// 	if !ok {
// 		fmt.Println("eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
// 	}
// 	b.FreeCell.headPtr.Status = Dead
// 	b.FreeCell.headPtr.SetIndex(FreeCellIndex)
// 	b.FreeCell.headPtr.Len = freeLen - newCell.headPtr.Len
// 	b.FreeCell.Tail.Tlen = b.FreeCell.headPtr.Len

// 	return offset
// }

// func (b *Bucket2) Get(offset uint16) (data []byte) {
// 	//bucket.Mutex.Lock()
// 	if offset == NILOFFSET {
// 		panic(fmt.Errorf("Bucket Select Offset Invalid"))
// 	}
// 	Len := *(*uint16)(unsafe.Pointer(&b.cellar[offset]))
// 	cell, ok := b.MakeCell(offset, Len-8-2)
// 	if ok && cell.headPtr.Len == cell.Tail.Tlen && cell.headPtr.Status {
// 		return unsafe.Slice(cell.bodyPtr, cell.headPtr.Len-8-2)
// 	} else {
// 		panic("rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr")
// 	}
// 	return nil
// 	//return s2b(b2s(body)) //unsafeCloneBytes(body): must be cloned before return
// 	//bucket.Mutex.UnLock()
// }

// func (bucket *Bucket2) CellHeader(offset uint32) *Head {
// 	return (*Head)(unsafe.Pointer(&bucket.cellar[offset:][0]))
// }

// func (bucket *Bucket2) IsRequiredToClean() bool {
// 	//return false
// 	return bucket.DeadSpace > BucketSize*5/100 //more than %5 waste should be defregmented
// }

// func (b *Bucket2) Offset(p Cell2) (offset uint16) {
// 	current := uintptr(unsafe.Pointer(p.headPtr))
// 	start := uintptr(unsafe.Pointer(&(b.cellar[0])))
// 	distance := int(current - start)
// 	if distance < 0 || distance > len(b.cellar) { //for debug purpose
// 		panic(fmt.Errorf("cell position is out of bucket buffer"))
// 	}
// 	return uint16(distance)
// }

// func (b *Bucket2) CheckLen() {
// 	for offset := 1; offset < len(b.cellar); {
// 		HLen := *(*uint16)(unsafe.Pointer(&b.cellar[offset]))
// 		TLen := *(*uint16)(unsafe.Pointer(&b.cellar[offset+int(HLen)-2]))
// 		if HLen != TLen {
// 			panic("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
// 		}
// 		offset += int(HLen)
// 	}
// 	for offset := len(b.cellar); offset > 1; {
// 		TLen := *(*uint16)(unsafe.Pointer(&b.cellar[offset-2]))
// 		HLen := *(*uint16)(unsafe.Pointer(&b.cellar[offset-int(TLen)]))
// 		if HLen != TLen {
// 			panic("-----------------------------------")
// 		}
// 		offset -= int(TLen)
// 	}
// }

// func (b *Bucket2) Delete(offset uint16) {
// 	Len := *(*uint16)(unsafe.Pointer(&b.cellar[offset]))
// 	cell, ok := b.MakeCell(offset, Len-8-2)
// 	if !ok || cell.headPtr.Len != cell.Tail.Tlen {
// 		panic("courrepted cell data")
// 	}

// 	cell.headPtr.Status = Dead
// 	b.DeadSpace += Len
// 	//return
// 	//mainOffset := offset

// 	freeLen := b.FreeCell.headPtr.Len
// 	if freeLen != b.FreeCell.Tail.Tlen {
// 		panic("freCell is invalid")
// 	}

// 	for next, ok := b.Next(cell); ok && !next.headPtr.Status; next, ok = b.Next(cell) {
// 		cell.MergeNext(next)
// 	}
// 	for prev, ok := b.Prev(cell); ok && !prev.headPtr.Status; prev, ok = b.Prev(cell) {
// 		prev.MergeNext(cell)
// 		cell = prev
// 	}
// 	if cell.headPtr.Len > freeLen {
// 		b.SetFreeCell(cell)
// 	}
// 	//cell.Tail.Len = cell.Head.Len
// }

// func (b *Bucket2) SetFreeCell(cell Cell2) {
// 	freeOffset := b.Offset(b.FreeCell)
// 	freeLen := b.FreeCell.headPtr.Len
// 	b.FreeCell.headPtr = cell.headPtr
// 	b.FreeCell.Tail = cell.Tail
// 	if b.FreeCell.Tail.Tlen != b.FreeCell.headPtr.Len {
// 		panic("cell is invalid")
// 	}
// 	b.DeadSpace -= cell.headPtr.Len - freeLen

// 	offset := b.Offset(cell)
// 	if freeOffset < offset || int(freeOffset) > int(offset)+int(cell.headPtr.Len) {
// 		b.DeadSpace += freeLen
// 	}
// }

// func (b *Bucket2) Next(cell Cell2) (next Cell2, ok bool) {
// 	nextOffset := b.Offset(cell) + cell.headPtr.Len //cell.Sizeof()
// 	if nextOffset > uint16(len(b.cellar)-10) {
// 		return next, false
// 	}
// 	nextLen := *(*uint16)(unsafe.Pointer(&b.cellar[nextOffset]))
// 	next, ok = b.MakeCell(nextOffset, nextLen-2-8)
// 	if next.headPtr.Len != next.Tail.Tlen {
// 		panic("next is invalid")
// 	}
// 	return
// }

// func (b *Bucket2) Prev(cell Cell2) (prev Cell2, ok bool) {
// 	offset := b.Offset(cell)
// 	if offset < 1+10 {
// 		return prev, false
// 	}
// 	prevLen := *(*uint16)(unsafe.Pointer(&b.cellar[offset-2]))
// 	prevOffset := offset - prevLen
// 	prev, ok = b.MakeCell(prevOffset, prevLen-2-8)
// 	if prev.headPtr.Len != prev.Tail.Tlen {
// 		panic("prev is invalid")
// 	}
// 	return
// }

// func (b *Bucket2) FindMaxDead() (Cell2, bool) {
// 	var Len, Max, maxOffset uint16
// 	var offset int
// 	for offset = 1; offset <= len(b.cellar)-10; offset = offset + int(Len) {
// 		Len = *(*uint16)(unsafe.Pointer(&b.cellar[offset]))
// 		Status := *(*bool)(unsafe.Pointer(&b.cellar[offset+2]))
// 		if !Status && Len > Max {
// 			Max = Len
// 			maxOffset = uint16(offset)
// 		}
// 	}
// 	if maxOffset > 0 {
// 		return b.MakeCell(maxOffset, Max-2-8)
// 	}
// 	return Cell2{}, false
// }

// func (b *Bucket2) Defrag(addr map[int64]uint16) {
// 	var Len uint16
// 	var offset int
// 	//tmpBucket.Reset()
// 	dst := b.cellar[1:]
// 	dst = dst[:0]
// 	for offset = 1; offset <= len(b.cellar)-10; offset = offset + int(Len) {
// 		Len = *(*uint16)(unsafe.Pointer(&b.cellar[offset]))
// 		Status := *(*bool)(unsafe.Pointer(&b.cellar[offset+2]))
// 		Index := int64(*(*uint32)(unsafe.Pointer(&b.cellar[offset+3])))
// 		Index |= int64(*(*byte)(unsafe.Pointer(&b.cellar[offset+7]))) << 32
// 		if Status {
// 			newOffset := 1 + len(dst)
// 			dst = append(dst, b.cellar[offset:offset+int(Len)]...)
// 			addr[Index] = uint16(newOffset)
// 		}
// 	}

// 	//tmpBucket.cellar, b.cellar = b.cellar, tmpBucket.cellar //swap

// 	freeLen := uint16(cap(b.cellar) - 1 - len(dst))
// 	b.FreeCell.headPtr = (*Head)(unsafe.Pointer(&b.cellar[1+len(dst)]))
// 	b.FreeCell.headPtr.Len = freeLen
// 	b.FreeCell.headPtr.Status = Dead
// 	b.FreeCell.headPtr.SetIndex(FreeCellIndex)
// 	b.FreeCell.bodyPtr = (*byte)(unsafe.Pointer(&b.cellar[1+len(dst)+8]))
// 	b.FreeCell.Tail = (*Tail)(unsafe.Pointer(&b.cellar[cap(b.cellar)-2]))
// 	b.FreeCell.Tail.Tlen = freeLen
// 	b.DeadSpace = 0
// }

// func (b *Bucket2) CheckAddr(addr map[int64]uint16) {
// 	for index, offset := range addr {
// 		if offset != NILOFFSET && string(b.Get(offset)) != fmt.Sprint(index) {
// 			panic("check read")
// 		}
// 	}
// }
