package basic

import (
	"fmt"

	"github.com/nazarifard/bigtype/internal/bucket"
	"github.com/nazarifard/bigtype/log"
)

type BitArray struct {
	buckets      [][]byte
	unexpandable bool
	Size         int
}

func NewBitArray(ops ...any) Array[bool] {
	ba := new(BitArray)
	opt := ParsArrayOptions[bool](ops...)
	ba.expand(opt.Size())
	ba.unexpandable = !opt.Expandable()
	return ba
}

func (ba *BitArray) Get(index int) (flag bool) {
	if index >= ba.Len() {
		panic(fmt.Errorf("index:%d is out of range. MaxSize of BitArray is %d", index, ba.Len()))
	}
	row := index / 8 / bucket.BucketSize
	column := (index / 8) % bucket.BucketSize
	r := index & 0x07 //  index%8
	return ba.buckets[row][column]&(1<<r) != 0
}

func (ba *BitArray) Set(index int, flag bool) {
	if index >= ba.Len() {
		panic(fmt.Errorf("index:%d is out of range. MaxSize of BitArray is %d", index, ba.Len()))
	}
	row := index / 8 / bucket.BucketSize
	column := index / 8 % bucket.BucketSize
	r := index & 0x07 //index%8
	if flag {
		ba.buckets[row][column] |= (1 << r)
	} else {
		ba.buckets[row][column] &^= (1 << r)
	}
}

func (ba *BitArray) Update(index int, updateFn func(old bool) (new bool)) {
	ba.Set(index, updateFn(ba.Get(index)))
}

func (ba *BitArray) Len() int {
	return ba.Size
}

func (ba *BitArray) expand(size int) {
	if size < ba.Size {
		return
	}

	if size < ba.Cap() {
		ba.Size = size
		return //there is enough empty space yet
	}

	//usually last bucket has empty space
	requierdBits := size - ba.Cap()
	requiredSpace := (requierdBits + 7) / 8
	for requiredSpace > 0 {
		row := make([]byte, bucket.BucketSize)
		ba.buckets = append(ba.buckets, row)
		requiredSpace -= bucket.BucketSize
		if log.VerboseMode {
			log.Logger.Info(fmt.Sprintf("new bucket allocated bucketSize=%d bytes, BucketLen:%d", bucket.BucketSize, bucket.BucketSize))
		}
	}
	ba.Size = size
}

func (ba *BitArray) Cap() int {
	return len(ba.buckets) * bucket.BucketSize * 8
}

func (ba *BitArray) Delete(index int) {
	ba.Set(index, false)
}
