package addr

import (
	"fmt"
	"unsafe"

	"github.com/nazarifard/bigtype/log"
)

const BucketSize = 64*1024 - 1             //16MB //Max:2*1024*1024*1024-1
const MaxAvailiableAddress = (1 << 40) - 1 //5 byte

type AddressItem struct {
	bucketNum [3]byte
	offset    [2]byte //uint16
}

func NewAddrItem(bucktetNum uint32, offset uint16) AddressItem {
	var ai AddressItem
	ai.bucketNum[0] = byte(bucktetNum)
	ai.bucketNum[1] = byte(bucktetNum >> 8)
	ai.bucketNum[2] = byte(bucktetNum >> 16)

	ai.offset[0] = byte(offset)
	ai.offset[1] = byte(offset >> 8)
	return ai
}

func (i AddressItem) BucketNum() uint32 {
	return uint32(i.bucketNum[2])<<16 |
		uint32(i.bucketNum[1])<<8 |
		uint32(i.bucketNum[0])
}

func (i AddressItem) Offset() uint16 {
	return uint16(i.offset[1])<<8 |
		uint16(i.offset[0])
}

type AddrBucket = [BucketSize / unsafe.Sizeof(AddressItem{})]AddressItem

type AddressTable struct {
	rows      []*AddrBucket
	FixedSize bool //Unexpandable
	Len       int
}

type IndexAddrPair struct {
	Index int
	Addr  AddressItem
}

func (at *AddressTable) Cap() int {
	return len(at.rows) * len(AddrBucket{})
	//bucket.BucketSize / int(unsafe.Sizeof(AddressItem{}))
}

func (at *AddressTable) Expand(size int) {
	if at.FixedSize {
		panic(fmt.Errorf("unexpandable AddressTable can not be expand"))
	}
	if size <= at.Cap() {
		at.Len = size
		return
	}

	var i int
	need := size / len(AddrBucket{})
	if size%(len(AddrBucket{})) > 0 {
		need++ // for remained
	}
	existing := len(at.rows)
	lack := need - existing
	newAddrPtr := make([]*AddrBucket, lack)
	at.rows = append(at.rows, newAddrPtr...)

	rows := make([]AddrBucket, lack)
	if log.VerboseMode {
		log.Logger.Info(fmt.Sprintf("taked %v bytes memory for a new address bucket", lack*int(unsafe.Sizeof(rows[0]))))
	}

	for i = 0; i < lack; i++ {
		at.rows[existing+i] = &rows[i]
	}

	at.Len = size
}

func (at *AddressTable) Expandable() bool {
	return !at.FixedSize
}

func (at *AddressTable) Set(index int, item AddressItem) {
	if !at.FixedSize && at.Len <= index {
		at.Expand(int(index) + 1)
	}
	row := index / len(AddrBucket{})
	column := index % len(AddrBucket{})
	at.rows[row][column] = item
}

func (at *AddressTable) Get(index int) AddressItem {
	row := index / len(AddrBucket{})    //uint32(len(at.Rows))
	column := index % len(AddrBucket{}) //uint32(len(at.Rows))
	item := at.rows[row][column]
	return item
}

func (at *AddressTable) UpdateIndexes(items []IndexAddrPair) {
	for i := range items {
		index := items[i].Index
		AddressItem := items[i].Addr
		at.Set(index, AddressItem)
	}
}
