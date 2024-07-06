package options

import (
	"fmt"
	"strconv"

	marshal "github.com/nazarifard/marshaltap"
	"github.com/nazarifard/marshaltap/tap/stdlib/gob"
)

type ArrayOptions[V any] struct {
	Size       int
	VMarshal   marshal.Interface[V]
	Expandable bool
}

func ParseArrayOptions[V any](params ...any) (o ArrayOptions[V]) {
	o.Size = 0                   //default
	o.VMarshal = gob.GobTap[V]{} //default
	o.Expandable = true

	switch len(params) {
	case 3:
		if expandable, ok := params[2].(bool); ok {
			o.Expandable = expandable
		} else {
			panic(fmt.Errorf("ParseArrayOptions failed"))
		}
		if size, ok := params[0].(int); ok {
			o.Size = size
			if params[1] == nil {
				return //Default
			}
			if m, ok := params[1].(marshal.Interface[V]); ok {
				o.VMarshal = m
				return
			}
		}

	case 2:
		var err error
		o.Size, err = parseSize(params[0])
		if err != nil {
			panic(fmt.Errorf("parsArrayOptions failed. cause: %w", err))
		}
		if o.Size > 0 {
			o.Expandable = false
		}
		if params[1] == nil {
			return //Default
		}
		if m, ok := params[1].(marshal.Interface[V]); ok {
			o.VMarshal = m
			return
		}

	case 1:
		switch t := params[0].(type) {
		default:
			var err error
			o.Size, err = parseSize(t)
			if err != nil {
				panic(fmt.Errorf("parsArrayOptions failed. cause: %w", err))
			}
			if o.Size > 0 {
				o.Expandable = false
			}
			return

		case marshal.Interface[V]:
			o.VMarshal = t
			return
		}
	case 0:
		return
	}
	panic(fmt.Errorf("parseArrayOptions failed. (size int)|(size int, m marshalInterface) are valid format"))
}

func parseSize(v any) (int, error) {
	s := fmt.Sprint(v)
	return strconv.Atoi(s)
}

type TreeOptions[K comparable, V any] struct {
	HistSize int
	VMarshal marshal.Interface[V]
}

func ParseTreeOptions[K comparable, V any](params ...any) (to TreeOptions[K, V]) {
	arrOption := ParseArrayOptions[V](params...)
	to.HistSize = arrOption.Size
	to.VMarshal = arrOption.VMarshal
	return
}

type MapOptions[K comparable, V any] struct {
	HintSize       int
	KMarshal       marshal.Interface[K]
	VMarshal       marshal.Interface[V]
	CheckCollition bool
}

func ParseMapOptions[K comparable, V any](params ...any) (mo MapOptions[K, V]) {
	mo.HintSize = 0
	//mo.KVMarshal = TapKV[]
	mo.CheckCollition = false

	var tk marshal.Interface[K] = gob.GobTap[K]{}
	var tv marshal.Interface[V] = gob.GobTap[V]{}

	switch len(params) {
	case 4:
		if b, ok := params[3].(bool); ok {
			mo.CheckCollition = b
		}
		fallthrough
	case 3:
		if m, ok := params[2].(marshal.Interface[V]); ok {
			tv = m
		}
		fallthrough
	case 1, 2:
		var n int
		if len(params) >= 2 {
			n = 2
		} else {
			n = 1
		}
		arOption := ParseArrayOptions[K](params[:n]...)
		mo.HintSize = arOption.Size
		tk = arOption.VMarshal
		fallthrough
	case 0:
		mo.KMarshal = tk
		mo.VMarshal = tv
		return mo
	}
	panic(fmt.Errorf("array parameters are invalid. (size int)|(size int, ki,vi marshalInterface, checkColliton bool) are valid format"))
}

// type MapItem[K comparable, V any] struct {
// 	Key K
// 	Value V
// }
