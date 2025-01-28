package basic

import (
	"fmt"
	"iter"
	"math/rand"

	"github.com/nazarifard/bigtype/internal/addr"
	"github.com/nazarifard/bigtype/internal/bucket"
	"github.com/nazarifard/bigtype/log"
)

const NILOFFSET = 0

type BytesArray struct {
	buckets      []bucket.Bucket
	addressTable addr.AddressTable
	//extandable   bool
}

func Ptr[V any](v V) *V { return &v }
func Val[V any](p *V) V { return *p }

func (ba *BytesArray) UnsafePtr(index int) []byte {
	if index > ba.Len() {
		panic(fmt.Errorf("index:%d is out of range. MaxSize of BigArray:%T is %d", index, ba, ba.Len()))
	}
	addrItem := ba.addressTable.Get(index)
	if addrItem.Offset() == NILOFFSET {
		return nil //here offset is not out of range
	}
	bs := ba.buckets[addrItem.BucketNum()].Get(addrItem.Offset())
	return bs
}

// read-only unsafe-mode; for manipulation or reuse later must be cloned
func (ba *BytesArray) Get(index int) []byte {
	return ba.UnsafePtr(index)
	//return append([]byte(nil), ba.UnsafePtr(index)...) //return copy of data
}

func (ba *BytesArray) Len() int {
	return ba.addressTable.Len
}

func (ba *BytesArray) IsDefragRequired(bucketId int) bool {
	return ba.buckets[bucketId].Deadspace() >= bucket.BucketSize*5/100
}

func (ba *BytesArray) DefragBucketCandidate(minNeededSpace int) (bucketId int) {
	for range len(ba.buckets) {
		bucketId = int(rand.Int31n(int32(len(ba.buckets))))
		if ba.IsDefragRequired(bucketId) &&
			ba.buckets[bucketId].Deadspace() >= minNeededSpace {
			return bucketId
		}
	}
	return -1 //all buckets are almost full
}

//	func (ba *BigArray[T]) GarbageCollector() {
//		for {
//			maxWastedIndex := ba.MaxWastedBucketNum()
//			if ba.Buckets[maxWastedIndex].isRequiredToClean() {
//				ba.clean(maxWastedIndex)
//			} else {
//				time.Sleep(200) //???
//			}
//		}
//	}

func (ba *BytesArray) addNewBucket() {
	b := bucket.NewBucket1(len(ba.buckets))
	ba.buckets = append(ba.buckets, b)

	if log.VerboseMode {
		sizeOf := bucket.BucketSize * len(ba.buckets)
		log.Logger.Info(fmt.Sprintf("BigArray taked a new Bucket, ArraySize:=%v NumOfBuckets:%d", sizeOf, len(ba.buckets)))
	}
}

func (ba *BytesArray) insertRequest(key uint32, Len int) (space []byte, ad addr.AddressItem) {
	if len(ba.buckets) == 0 {
		ba.addNewBucket()
		last := len(ba.buckets) - 1
		space, offset := ba.buckets[last].Request(int(key), Len)
		return space, addr.NewAddrItem(uint32(last), offset)
	}
	//first try last bucket
	i := len(ba.buckets) - 1
	for range len(ba.buckets) {
		if ba.buckets[i] != nil {
			space, offset := ba.buckets[i].Request(int(key), Len)
			if offset != NILOFFSET {
				return space, addr.NewAddrItem(uint32(i), offset)
			}
		}
		//exclude last bucket; its checked firstly
		if len(ba.buckets) > 1 {
			i = int(rand.Int31n(int32(len(ba.buckets)) - 1)) //try for next bucket randomly
		}
	}

	//all buckets is nil or full
	cellInfoSize := ba.buckets[0].InfoSize()
	minimumRequired := cellInfoSize + Len
	defragId := ba.DefragBucketCandidate(minimumRequired)
	if 0 <= defragId && defragId < len(ba.buckets) {
		ba.buckets[defragId].Defrag(&ba.addressTable)
		space, offset := ba.buckets[defragId].Request(int(key), Len)
		if offset == NILOFFSET {
			panic("defrag problem")
		}
		return space, addr.NewAddrItem(uint32(defragId), offset)
	} else {
		ba.addNewBucket()
		last := len(ba.buckets) - 1
		space, offset := ba.buckets[last].Request(int(key), Len)
		return space, addr.NewAddrItem(uint32(last), offset)
	}
}

// func (ba *BytesArray) insert(key uint32, value []byte) addr.AddressItem {
// 	if len(ba.buckets) == 0 {
// 		ba.addNewBucket()
// 		last := len(ba.buckets) - 1
// 		offset := ba.buckets[last].Write(int(key), value)
// 		return addr.NewAddrItem(uint32(last), offset)
// 	}
// 	//first try last bucket
// 	i := len(ba.buckets) - 1
// 	for range len(ba.buckets) {
// 		if ba.buckets[i] != nil {
// 			offset := ba.buckets[i].Write(int(key), value)
// 			if offset != NILOFFSET {
// 				return addr.NewAddrItem(uint32(i), offset)
// 			}
// 		}
// 		//exclude last bucket; its checked firstly
// 		if len(ba.buckets) > 1 {
// 			i = int(rand.Int31n(int32(len(ba.buckets)) - 1)) //try for next bucket randomly
// 		}
// 	}

// 	//all buckets is nil or full
// 	cellInfoSize := ba.buckets[0].InfoSize()
// 	minimumRequired := cellInfoSize + len(value)
// 	defragId := ba.DefragBucketCandidate(minimumRequired)
// 	if 0 <= defragId && defragId < len(ba.buckets) {
// 		ba.buckets[defragId].Defrag(&ba.addressTable)
// 		offset := ba.buckets[defragId].Write(int(key), value)
// 		if offset == NILOFFSET {
// 			panic("defrag problem")
// 		}
// 		return addr.NewAddrItem(uint32(defragId), offset)
// 	} else {
// 		ba.addNewBucket()
// 		last := len(ba.buckets) - 1
// 		offset := ba.buckets[last].Write(int(key), value)
// 		return addr.NewAddrItem(uint32(last), offset)
// 	}
// }

func (ba *BytesArray) Delete(index int) {
	num := ba.addressTable.Get(index).BucketNum()
	b := ba.buckets[num]
	b.Delete(ba.addressTable.Get(index).Offset())  //delete from bucket
	ba.addressTable.Set(index, addr.AddressItem{}) //delete from address table
	// if bucket.isRequiredToClean() {
	// 	ba.clean(num)
	// }
}

//	func (ba *BigArray[T]) checkUp() {
//		ba.operationCounter++
//		if ba.operationCounter == 1000*1000 {
//			runtime.GC()
//			ba.operationCounter = 0
//		}
//	}
func (ba *BytesArray) Request(index int, Len int) []byte {
	if ba.addressTable.Expandable() && index >= ba.Len() {
		ba.addressTable.Expand(int(index) + 1)
	} else if index >= ba.Len() {
		panic(fmt.Errorf("index:%d is out of range. MaxSize of BigArray:%T is %d", index, ba, ba.Len()))
	}
	addrItem := ba.addressTable.Get(index)
	var oldValue []byte
	if addrItem.Offset() != NILOFFSET {
		oldValue = ba.buckets[addrItem.BucketNum()].Get(addrItem.Offset())
	}
	switch len(oldValue) {
	case 0: //insert new record
		space, addrItem := ba.insertRequest(uint32(index), Len)
		ba.addressTable.Set(index, addrItem)
		return space

	case Len: //just replace
		return oldValue //copy(oldValue, value)

	default: //delete + insert
		ba.Delete(index)
		space, addrItem := ba.insertRequest(uint32(index), Len)
		ba.addressTable.Set(index, addrItem)
		return space
	}
}
func (ba *BytesArray) Set(index int, value []byte) {
	space := ba.Request(index, len(value))
	copy(space, value)
}

func (ba *BytesArray) Update(index int, fn func(old []byte) (new []byte)) {
	if index > ba.Len() {
		panic(fmt.Errorf("index:%d is out of range. MaxSize of BigArray:%T is %d", index, ba, ba.Len()))
	}
	var v []byte
	addrItem := ba.addressTable.Get(index)
	if addrItem.Offset() == NILOFFSET {
		ba.Set(index, fn(v))
	}
	bs := ba.buckets[addrItem.BucketNum()].Get(addrItem.Offset())
	ba.Set(index, fn(bs))
}

func NewBytesArray(size int, extandable bool) *BytesArray { //Array[[]byte] {
	var ba BytesArray
	ba.addressTable.Expand(size)
	ba.addressTable.FixedSize = !extandable
	return &ba
}

func (ba *BytesArray) All() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for i := range ba.Len() {
			if !yield(ba.UnsafePtr(i)) {
				break
			}
		}
	}
}
