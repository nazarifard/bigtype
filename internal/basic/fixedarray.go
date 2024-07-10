package basic

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/nazarifard/bigtype/internal/bucket"
	"github.com/nazarifard/bigtype/log"
)

type BigFixedArray[V any] struct {
	buckets      [][]V
	unexpandable bool
	Size         int
}

func (ba *BigFixedArray[V]) Cap() int {
	var v V
	columns := bucket.BucketSize / int(unsafe.Sizeof(v))
	rows := len(ba.buckets)
	return rows * columns
}

func NewFixedArray[V any](size int, extendable bool) Array[V] {
	if !IsFixedType(*new(V)) {
		panic(fmt.Errorf("type:%T is not fixed sized type. this type of bigarray doesn't support dynamic size arrays", *new(V)))
	}
	if reflect.ValueOf(*new(V)).Kind() == reflect.Bool {
		var options ArrayOptions[V]
		options.WithSize(size).WithExtandable(extendable)
		return NewBitArray(options).(Array[V])
	}
	var ba BigFixedArray[V]
	ba.expand(size)
	ba.unexpandable = !extendable ////root is reserved in tree use case
	return &ba
}

func (ba *BigFixedArray[V]) Sizeof() int {
	items := 0
	rows := len(ba.buckets)
	for i := range rows {
		items += len(ba.buckets[i])
	}
	return items * int(unsafe.Sizeof(*new(V)))
}

func (ba *BigFixedArray[V]) expand(size int) {
	if size < ba.Size {
		return
	}

	if size < ba.Cap() {
		ba.Size = size
		return //there is enough empty space yet
	}

	//usually last bucket has empty space
	itemSize := int(unsafe.Sizeof(*new(V)))
	requiredSpace := itemSize * (size - ba.Cap())
	for requiredSpace > -bucket.BucketSize {
		columns := bucket.BucketSize / int(unsafe.Sizeof(*new(V)))
		row := make([]V, columns)
		ba.buckets = append(ba.buckets, row)
		requiredSpace -= bucket.BucketSize
		if log.VerboseMode {
			log.Logger.Info(fmt.Sprintf("new bucket allocated bucketSize=%d bytes, BucketLen:%d", bucket.BucketSize, columns))
		}
	}
	ba.Size = size
}

func (ba *BigFixedArray[V]) Len() int {
	return int(ba.Size)
}

func (ba *BigFixedArray[V]) Set(index int, item V) {
	if !ba.unexpandable && index >= ba.Size {
		ba.expand(index + 1)
	}

	if len(ba.buckets) == 0 {
		panic(fmt.Errorf("bigArray is empty"))
	}
	if index > ba.Size {
		panic(fmt.Errorf("index:%d is out of range. legal size:%d", index, ba.Size))
	}
	bucketSize := len(ba.buckets[0]) //BucketSize := BucketSize + (BucketSize % unsafe.Sizeof(item))
	row := index / bucketSize
	column := index % bucketSize
	if int(row) >= len(ba.buckets) || int(column) >= len(ba.buckets[row]) {
		fmt.Println("bucketSize, index, row, column", bucketSize, index, row, column)
	}

	ba.buckets[row][column] = item
}

func (ba *BigFixedArray[V]) Get(index int) V {
	if len(ba.buckets) == 0 {
		panic(fmt.Errorf("bigArray is empty"))
	}
	if index > ba.Size {
		panic(fmt.Errorf("index:%d is out of range. legal size:%d", index, ba.Size))
	}
	bucketSize := len(ba.buckets[0])
	row := index / bucketSize
	column := index % bucketSize
	return ba.buckets[row][column]
}

func (ba *BigFixedArray[V]) Update(index int, updateFn func(old V) (new V)) {
	ba.Set(index, updateFn(ba.Get(index)))
}

func (ba *BigFixedArray[V]) Delete(index int) {
	ba.Set(index, *new(V))
}

// BigFixedArray does not support UnsafePtr and does not need at all
func (ba *BigFixedArray[V]) UnsafePtr(index int) *V {
	return nil
}
