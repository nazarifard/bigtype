package basic

import (
	"fmt"
	"math/rand"
	"unsafe"

	"github.com/nazarifard/bigtype/internal/addr"
	"github.com/nazarifard/bigtype/internal/bucket"
	"github.com/nazarifard/bigtype/internal/options"
	"github.com/nazarifard/bigtype/internal/utils"
	"github.com/nazarifard/bigtype/log"
	marshal "github.com/nazarifard/marshaltap"
)

const NILOFFSET = 0

type array[V any] struct {
	marshaller    marshal.Interface[V]
	buckets       []bucket.Bucket
	addressTable  addr.AddressTable
	cap           int
	isStringArray bool //V ~string
	isBytesArray  bool //V ~[]~byte
	//tempBucket    *bucket.Bucket
}

func (ba *array[V]) Len() int {
	return ba.addressTable.Len
}

func (ba *array[V]) IsDefragRequired(bucketId int) bool {
	return ba.buckets[bucketId].Deadspace() >= bucket.BucketSize*5/100
}

func (ba *array[V]) DefragBucketCandidate(minNeededSpace int) (bucketId int) {
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

func (ba *array[V]) addNewBucket() {
	b := bucket.NewBucket1(len(ba.buckets))
	ba.buckets = append(ba.buckets, b)

	if log.VerboseMode {
		sizeOf := bucket.BucketSize * len(ba.buckets)
		log.Logger.Info(fmt.Sprintf("BigArray taked a new Bucket, ArraySize:=%v NumOfBuckets:%d", sizeOf, len(ba.buckets)))
	}
}

func (ba *array[V]) insert(key uint32, value []byte) addr.AddressItem {
	if len(ba.buckets) == 0 {
		ba.addNewBucket()
		last := len(ba.buckets) - 1
		offset := ba.buckets[last].Write(int(key), value)
		return addr.NewAddrItem(uint32(last), offset)
	}
	//first try last bucket
	i := len(ba.buckets) - 1
	for range len(ba.buckets) {
		if ba.buckets[i] != nil {
			offset := ba.buckets[i].Write(int(key), value)
			if offset != NILOFFSET {
				return addr.NewAddrItem(uint32(i), offset)
			}
		}
		//exclude last bucket; its checked firstly
		if len(ba.buckets) > 1 {
			i = int(rand.Int31n(int32(len(ba.buckets)) - 1)) //try for next bucket randomly
		}
	}

	//all buckets is nil or full
	cellInfoSize := ba.buckets[0].InfoSize()
	minimumRequired := cellInfoSize + len(value)
	defragId := ba.DefragBucketCandidate(minimumRequired)
	if 0 <= defragId && defragId < len(ba.buckets) {
		ba.buckets[defragId].Defrag(&ba.addressTable)
		offset := ba.buckets[defragId].Write(int(key), value)
		if offset == NILOFFSET {
			panic("defrag problem")
		}
		return addr.NewAddrItem(uint32(defragId), offset)
	} else {
		ba.addNewBucket()
		last := len(ba.buckets) - 1
		offset := ba.buckets[last].Write(int(key), value)
		return addr.NewAddrItem(uint32(last), offset)
	}
}

func (ba *array[V]) delete(index int) {
	num := ba.addressTable.Get(index).BucketNum()
	b := ba.buckets[num]
	b.Delete(ba.addressTable.Get(index).Offset())  //delete from bucket
	ba.addressTable.Set(index, addr.AddressItem{}) //delete from address table
	// if bucket.isRequiredToClean() {
	// 	ba.clean(num)
	// }
}

// func (ba *BigArray[T]) checkUp() {
// 	ba.operationCounter++
// 	if ba.operationCounter == 1000*1000 {
// 		runtime.GC()
// 		ba.operationCounter = 0
// 	}
// }

func (ba *array[V]) Set(index int, v V) {
	if ba.addressTable.Expandable() && index >= ba.Len() {
		ba.addressTable.Expand(int(index) + 1)
	} else if index >= ba.Len() {
		panic(fmt.Errorf("index:%d is out of range. MaxSize of BigArray:%T is %d", index, ba, ba.cap))
	}

	var value []byte
	if ba.isStringArray {
		s := *(*string)(unsafe.Pointer(&v))
		value = s2b(s) //here s2b is correct
	} else if ba.isBytesArray {
		value = *(*[]byte)(unsafe.Pointer(&v))
	} else {
		wbuf, err := ba.marshaller.Encode(v)
		if err != nil {
			panic(fmt.Errorf("ba.marshaller.Encode Faild. Cause: %w", err))
		}

		defer wbuf.Free()
		// if bucket.BucketSize < wbuf.Len()+int(bucket.CellHeaderSize) {
		// 	panic(fmt.Errorf("too big object. type:%T data:%v len(v)=%d must be smaller than %d", v, wbuf, wbuf.Len(), bucket.BucketSize-bucket.CellHeaderSize))
		// }
		value = wbuf.Bytes()
	}

	// if bucket.BucketSize < len(value)+int(bucket.CellHeaderSize) {
	// 	panic(fmt.Errorf("too big object. type:%T data:%v len(v)=%d must be smaller than %d", v, v, len(value), bucket.BucketSize-bucket.CellHeaderSize))
	// }

	addrItem := ba.addressTable.Get(index)
	var oldValue []byte
	if addrItem.Offset() != NILOFFSET {
		oldValue = ba.buckets[addrItem.BucketNum()].Read(addrItem.Offset())
	}
	switch len(oldValue) {
	case 0: //insert new record
		addrItem := ba.insert(uint32(index), value)
		ba.addressTable.Set(index, addrItem)

	case len(value): //just replace
		copy(oldValue, value)

	default: //delete + insert
		ba.delete(index)
		addrItem := ba.insert(uint32(index), value)
		ba.addressTable.Set(index, addrItem)
	}
}

func (ba *array[V]) Update(index int, updateFn func(old V) (new V)) {
	if index > ba.Len() {
		panic(fmt.Errorf("index:%d is out of range. MaxSize of BigArray:%T is %d", index, ba, ba.Len()))
	}
	var v V
	addrItem := ba.addressTable.Get(index)
	if addrItem.Offset() == NILOFFSET {
		ba.Set(index, updateFn(v))
	}
	bs := ba.buckets[addrItem.BucketNum()].Read(addrItem.Offset())
	if ba.isStringArray {
		//skip decode
		pv := (*string)(unsafe.Pointer(&v))
		*pv = unsafe.String(&bs[0], len(bs)) //yes unsafe.String(&bs[0], len(bs))
	} else if ba.isBytesArray {
		//skip decode
		v = *(*V)(unsafe.Pointer(&bs))
	} else {
		var err error
		v, _, err = ba.marshaller.Decode(bs)
		if err != nil {
			panic(fmt.Sprintf("Unmarshal Error. Cause:%s", err.Error()))
		}
	}
	ba.Set(index, updateFn(v))
}

func (ba *array[V]) Get(index int) V {
	if index > ba.Len() {
		panic(fmt.Errorf("index:%d is out of range. MaxSize of BigArray:%T is %d", index, ba, ba.Len()))
	}

	addrItem := ba.addressTable.Get(index)
	if addrItem.Offset() == NILOFFSET {
		var t V
		return t //here offset is not out of range
	}
	bs := ba.buckets[addrItem.BucketNum()].Read(addrItem.Offset())

	var v V
	if ba.isStringArray {
		//skipm decode
		pv := (*string)(unsafe.Pointer(&v))
		*pv = string(bs) //not unsafe.String(&bs[0], len(bs))
	} else if ba.isBytesArray {
		//skip decode
		dst := make([]byte, len(bs))
		copy(dst, bs) //must be copied
		v = *(*V)(unsafe.Pointer(&dst))
	} else {
		var err error
		v, _, err = ba.marshaller.Decode(bs)
		if err != nil {
			panic(fmt.Sprintf("Unmarshal Error. Cause:%s", err.Error()))
		}
	}
	return v
}

func NewArray[V any](ops ...any) Array[V] {
	var v V
	if IsFixedType(v) {
		return NewFixedArray[V](ops...)
	}

	ba := array[V]{}
	option := options.ParseArrayOptions[V](ops...)
	ba.addressTable.Expand(option.Size)
	ba.addressTable.FixedSize = !option.Expandable
	//ba.tempBucket = bucket.NewBucket()
	ba.marshaller = option.VMarshal
	ba.isBytesArray = utils.IsBytes(v)
	ba.isStringArray = utils.IsString(v)

	return &ba
}
