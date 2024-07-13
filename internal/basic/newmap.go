package basic

import (
	"github.com/nazarifard/bigtype/internal/hash"
	"github.com/nazarifard/bigtype/internal/kv"
	"github.com/nazarifard/bigtype/internal/utils"
	"github.com/nazarifard/bigtype/options"
	"github.com/nazarifard/fastape"
)

func NewMap[K comparable, V any](ops ...any) Map[K, V] {
	opt := ParsMapOptions[K, V](ops...)
	if isNumber[K]() {
		var to options.TreeOptions[K, V]
		to.WithSize(opt.Size()).WithMarshal(opt.VTape())
		return makeTree[K, V](to)
	}

	if utils.IsFixedType(*new(K)) || utils.IsFixedType(*new(V)) {
		//bigmap1
		var vo options.ArrayOptions[V]
		var ko options.ArrayOptions[K]
		var ho options.TreeOptions[uint64, int]
		vo.WithSize(opt.Size()).WithMarshal(opt.VTape())
		ko.WithSize(opt.Size()).WithMarshal(opt.KTape())
		ho.WithSize(opt.Size())
		return &bigMap1[K, V]{
			htree:          newTree[uint64, int](ho), //TODO hint size
			collitionMap:   make(map[K]V),
			kHash:          hash.NewHash[K](nil),
			keys:           NewArray[K](ko),
			values:         NewArray[V](vo),
			checkCollition: opt.CollisionCheck(),
		}
	}

	//bigmap2
	var kvo options.ArrayOptions[kv.KV[K, V]]
	var ko options.ArrayOptions[K]
	var ho options.TreeOptions[uint64, int]
	kvo.WithSize(opt.Size()).WithMarshal(kv.NewTapeKV[K, V](opt.KTape(), opt.VTape()))
	ko.WithSize(opt.Size())
	ho.WithSize(opt.Size())

	if opt.VTape() == nil {
		if utils.IsBytes(*new(V)) {
			bsTape := any(fastape.SliceTape[byte, fastape.UnitTape[byte]]{}).(fastape.Tape[V])
			opt.WithVTape(bsTape)
		} else if utils.IsString(*new(V)) {
			strTape := any(fastape.StringTape{}).(fastape.Tape[V])
			opt.WithVTape(strTape)
		}
	}

	return &bigMap2[K, V]{
		htree:          newTree[uint64, int](ho), //TODO hint size
		collitionMap:   make(map[K]V),
		kHash:          hash.NewHash[K](nil),
		kv:             NewArray[kv.KV[K, V]](kvo),
		checkCollition: opt.CollisionCheck(),
		tape:           kv.NewTapeKV[K, V](opt.KTape(), opt.VTape()),
		vTape:          opt.VTape(),
	}
}
