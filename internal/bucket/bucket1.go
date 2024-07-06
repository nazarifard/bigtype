package bucket

import (
	"fmt"
	"unsafe"

	"github.com/nazarifard/bigtype/internal/addr"
	"github.com/nazarifard/bigtype/log"
)

const BucketSize = 64*1024 - 1 //uint16  //64KB //16MB //Max:2*1024*1024*1024-1
const NILOFFSET = 0
const FreeCellIndex = (1 << 40) - 1 // = 0xFFFFFFFFFF 5 bytes
const MaxValidIndex = FreeCellIndex - 1
const HEAD_SIZE = int(unsafe.Sizeof(*new(Head)))
const MaxValidOffset = BucketSize - 1

var _ = map[bool]bool{false: false, BucketSize <= 0xFFFF: true}
var _ = map[bool]bool{false: false, addr.BucketSize == BucketSize: true}

type BucketHeader1 struct {
	id        int //3 byte
	deadSpace uint16
	FreeCell  Cell1
	Len       int
}

type Bucket1 struct { //type Bucket1 [BucketSize]byte
	BucketHeader1
	cellar [BucketSize]byte
}

//var tmpBucket = Bucket1{}

func (b *Bucket1) MakeCell(offset, Len uint16) (c Cell1, ok bool) {
	if offset == NILOFFSET ||
		offset >= MaxValidOffset ||
		int(offset)+int(Len)+HEAD_SIZE > BucketSize {
		fmt.Println("offset", offset, "len", Len)
		return c, false
	}
	c.Head = (*Head)(unsafe.Pointer(&b.cellar[offset])) //cellar[0] is forbiden
	offset += uint16(unsafe.Sizeof(*c.Head))
	//last
	if offset > MaxValidOffset {
		c.Body = nil
	} else {
		c.Body = &b.cellar[offset]
	}
	//c.Tail = (*Tail)(unsafe.Pointer(&b.cellar[offset+Len]))
	return c, true
}
func (b *Bucket1) Id() int {
	return b.id
}
func (b *Bucket1) Reset() {
	b.BucketHeader1 = BucketHeader1{}
	b.FreeCell, _ = b.MakeCell(1, uint16(len(b.cellar)-1-HEAD_SIZE))
	b.FreeCell.Status = Dead
	b.FreeCell.SetIndex(FreeCellIndex)
	b.FreeCell.Head.Len = uint16(len(b.cellar) - 1)
	//b.FreeCell.Tail.Tlen = b.FreeCell.Head.Len
}

func NewBucket1(id int) *Bucket1 {
	if id > (1<<24)-1 {
		panic("bucket id is too big. must be smaller than 2^24")
	}
	bucket := new(Bucket1)
	bucket.Reset()
	bucket.id = id
	//fmt.Printf("\ntaked %v bytes memory for a new data bucket\n", cap(bucket.cellar))
	// if log.VerboseMode {
	// 	log.Logger.Info(fmt.Sprintf("taked %v bytes memory for a new data bucket", unsafe.Sizeof(*bucket)))
	// }
	return bucket
}

//	func (bucket *Bucket1) Bytes() []byte {
//		return bucket.cellar[:]
//	}
//
//	func (b *Bucket1) Set(key int, value []byte) (offset uint16) {
//		return b.write(key, value)
//	}
func (b *Bucket1) Write(index int, data []byte) (offset uint16) {
	freeOffset := b.Offset(b.FreeCell)
	freeLen := b.FreeCell.Head.Len

	if HEAD_SIZE+len(data) > int(freeLen)-HEAD_SIZE {
		return NILOFFSET
	}
	if index > MaxValidIndex {
		panic("index is too large more than 5 bytes.")
	}

	//b.FreeCell.Body = &b.cellar[freeOffset+uint16(HEAD_SIZE)] // : freeOffset+8-10+freeLen]

	newCell, _ := b.MakeCell(freeOffset, uint16(len(data)))
	newCell.Head.SetIndex(index)
	newCell.Head.Status = Live
	bs := unsafe.Slice(newCell.Body, len(data))
	copy(bs, data)
	newCell.Head.Len = uint16(HEAD_SIZE + len(bs))
	//newCell.Tail.Tlen = newCell.Head.Len
	offset = freeOffset //b.Offset(newCell)

	nextOffset := freeOffset + newCell.Head.Len
	var ok bool
	b.FreeCell, ok = b.MakeCell(nextOffset, freeLen-newCell.Head.Len-uint16(HEAD_SIZE))
	if !ok {
		fmt.Println("eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
	}
	b.FreeCell.Head.Status = Dead
	b.FreeCell.Head.SetIndex(FreeCellIndex)
	b.FreeCell.Head.Len = freeLen - newCell.Head.Len
	//b.FreeCell.Tail.Tlen = b.FreeCell.Head.Len

	return offset
}

func (b *Bucket1) Read(offset uint16) (data []byte) {
	//bucket.Mutex.Lock()
	if offset == NILOFFSET {
		panic(fmt.Errorf("Bucket1 Select Offset Invalid"))
	}
	Len := *(*uint16)(unsafe.Pointer(&b.cellar[offset]))
	cell, ok := b.MakeCell(offset, Len-uint16(HEAD_SIZE))
	if ok &&
		//cell.Head.Len == cell.Tail.Tlen &&
		cell.Head.Status {
		return unsafe.Slice(cell.Body, cell.Head.Len-uint16(HEAD_SIZE))
	} else {
		panic("rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr")
	}
	return nil
	//return s2b(b2s(body)) //unsafeCloneBytes(body): must be cloned before return
	//bucket.Mutex.UnLock()
}

func (bucket *Bucket1) CellHeader(offset uint32) *Head {
	return (*Head)(unsafe.Pointer(&bucket.cellar[offset:][0]))
}

func (bucket *Bucket1) IsRequiredToClean() bool {
	//return false
	return bucket.deadSpace > BucketSize*5/100 //more than %5 waste should be defregmented
}

func (b *Bucket1) Offset(p Cell1) (offset uint16) {
	current := uintptr(unsafe.Pointer(p.Head))
	start := uintptr(unsafe.Pointer(&(b.cellar[0])))
	distance := int(current - start)
	if distance < 0 || distance > len(b.cellar) { //for debug purpose
		panic(fmt.Errorf("cell position is out of bucket buffer"))
	}
	return uint16(distance)
}

func (b *Bucket1) CheckUp(ad addr.AddressTable) error {
	var offset, Len int
	for offset = 1; offset < cap(b.cellar); offset += Len {
		Len = int(*(*uint16)(unsafe.Pointer(&b.cellar[offset])))
		cell, ok := b.MakeCell(uint16(offset), uint16(Len-HEAD_SIZE))
		if !ok {
			return fmt.Errorf("checkup bucketId:%d failed. MakeCell return false", b.id)
		} else if offset == int(b.Offset(b.FreeCell)) {
			if cell.Head != b.FreeCell.Head {
				return fmt.Errorf("checkup bucketId:%d failed. FreeCell is incorrect", b.id)
			}
		} else if cell.Status && ad.Get(cell.Index()) != addr.NewAddrItem(uint32(b.id), uint16(offset)) {
			msg := fmt.Sprintf("checkup bucketId:%d failed. Index:%d at Offset:%d", b.id, cell.Index(), offset)
			msg += fmt.Sprintf("is not matched with addressTable.Offset:%d", ad.Get(cell.Index()))
			return fmt.Errorf(msg)
		}
	}
	if offset != cap(b.cellar) {
		return fmt.Errorf("checkup bucketId:%d failed. last_offset + Len != BucketSize", b.id)
	}
	return nil
}

func (b *Bucket1) Delete(offset uint16) {
	Len := *(*uint16)(unsafe.Pointer(&b.cellar[offset]))
	cell, ok := b.MakeCell(offset, Len-uint16(HEAD_SIZE))
	if !ok { // || cell.Head.Len != cell.Tail.Tlen {
		panic("courrepted cell data")
	}

	cell.Head.Status = Dead
	b.deadSpace += Len
	//return
	//mainOffset := offset

	freeLen := b.FreeCell.Head.Len
	// if freeLen != b.FreeCell.Tail.Tlen {
	// 	panic("freCell is invalid")
	// }

	for next, ok := b.Next(cell); ok && !next.Status; next, ok = b.Next(cell) {
		cell.MergeNext(next)
	}
	// for prev, ok := b.Prev(cell); ok && !prev.Status; prev, ok = b.Prev(cell) {
	// 	prev.MergeNext(cell)
	// 	cell = prev
	// }
	if cell.Head.Len > freeLen {
		b.SetFreeCell(cell)
	}
	//cell.Tail.Len = cell.Head.Len
}

func (b *Bucket1) SetFreeCell(cell Cell1) {
	freeOffset := b.Offset(b.FreeCell)
	freeLen := b.FreeCell.Head.Len
	b.FreeCell.Head = cell.Head
	//b.FreeCell.Tail = cell.Tail
	// if b.FreeCell.Tail.Tlen != b.FreeCell.Head.Len {
	// 	panic("cell is invalid")
	// }
	b.deadSpace -= cell.Head.Len - freeLen

	offset := b.Offset(cell)
	if freeOffset < offset || int(freeOffset) > int(offset)+int(cell.Head.Len) {
		b.deadSpace += freeLen
	}
}

func (b *Bucket1) Next(cell Cell1) (next Cell1, ok bool) {
	nextOffset := b.Offset(cell) + cell.Head.Len //cell.Sizeof()
	if nextOffset > uint16(len(b.cellar)-HEAD_SIZE) {
		return next, false
	}
	nextLen := *(*uint16)(unsafe.Pointer(&b.cellar[nextOffset]))
	next, ok = b.MakeCell(nextOffset, nextLen-uint16(HEAD_SIZE))
	// if next.Head.Len != next.Tail.Tlen {
	// 	panic("next is invalid")
	// }
	return
}

// func (b *Bucket1) Prev(cell Cell1) (prev Cell1, ok bool) {
// 	offset := b.Offset(cell)
// 	if offset < 1+10 {
// 		return prev, false
// 	}
// 	prevLen := *(*uint16)(unsafe.Pointer(&b.cellar[offset-2]))
// 	prevOffset := offset - prevLen
// 	prev, ok = b.MakeCell(prevOffset, prevLen-2-8)
// 	if prev.Head.Len != prev.Tail.Tlen {
// 		panic("prev is invalid")
// 	}
// 	return
// }

func (b *Bucket1) FindMaxDead() (Cell1, bool) {
	var Len, Max, maxOffset uint16
	var offset int
	for offset = 1; offset <= len(b.cellar)-HEAD_SIZE; offset = offset + int(Len) {
		Len = *(*uint16)(unsafe.Pointer(&b.cellar[offset]))
		Status := *(*bool)(unsafe.Pointer(&b.cellar[offset+2]))
		if !Status && Len > Max {
			Max = Len
			maxOffset = uint16(offset)
		}
	}
	if maxOffset > 0 {
		return b.MakeCell(maxOffset, Max-uint16(HEAD_SIZE))
	}
	return Cell1{}, false
}

func (b *Bucket1) Bytes() []byte {
	return b.cellar[:]
}

func (b *Bucket1) Defrag(ad *addr.AddressTable) {
	var Len uint16
	var offset int
	//tmpBucket.Reset()
	src := b.cellar
	dst := b.cellar[1:]
	dst = dst[:0]
	for offset = 1; offset <= len(src)-HEAD_SIZE; offset = offset + int(Len) {
		Len = *(*uint16)(unsafe.Pointer(&src[offset]))
		Status := *(*bool)(unsafe.Pointer(&src[offset+2]))
		Index := int(*(*uint32)(unsafe.Pointer(&src[offset+3])))
		Index |= int(*(*byte)(unsafe.Pointer(&src[offset+7]))) << 32
		if Status {
			newOffset := 1 + len(dst)
			dst = append(dst, src[offset:offset+int(Len)]...)
			oldOffset := ad.Get(Index)
			ad.Set(Index, addr.NewAddrItem(uint32(b.id), uint16(newOffset)))
			updatedOffset := ad.Get(Index)
			if oldOffset != updatedOffset {
				_ = updatedOffset
			}
		}
	}

	//tmpBucket.cellar, cellar = cellar, tmpBucket.cellar //swap
	freeLen := uint16(cap(src) - 1 - len(dst))
	freeOffset := 1 + len(dst)
	b.FreeCell.Head = (*Head)(unsafe.Pointer(&b.cellar[freeOffset]))
	b.FreeCell.Head.Len = freeLen
	b.FreeCell.Head.Status = Dead
	b.FreeCell.Head.SetIndex(FreeCellIndex)
	b.FreeCell.Body = nil
	if freeOffset+HEAD_SIZE <= MaxValidOffset {
		b.FreeCell.Body = (*byte)(unsafe.Pointer(&b.cellar[freeOffset+HEAD_SIZE]))
	}
	//b.FreeCell.Tail = (*Tail)(unsafe.Pointer(&cellar[cap(cellar)-2]))
	//b.FreeCell.Tail.Tlen = freeLen
	b.deadSpace = 0

	if log.VerboseMode {
		log.Logger.Info(fmt.Sprintf("defrag bucketId %d, Saved Space:%d %d%%", b.id, b.FreeCell.Len, b.FreeCell.Len*100/BucketSize))
	}

	//For Debug purpose
	// if err := b.CheckUp(*ad); err != nil {
	// 	panic(err)
	// }
}

func (b *Bucket1) Deadspace() int {
	return int(b.deadSpace)
}
func (b *Bucket1) InfoSize() int {
	return int(unsafe.Sizeof(*new(Head)))
}
