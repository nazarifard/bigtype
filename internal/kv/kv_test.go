package kv

import (
	"testing"
	"time"

	"github.com/nazarifard/fastape"
)

func TestKVss(t *testing.T) {
	sstape := NewTapeKV[string, string](nil, nil)
	item := KV[string, string]{
		Key:   "key",
		Value: "Value",
	}
	bs := make([]byte, sstape.Sizeof(item))
	n, err := sstape.Roll(item, bs)
	_, _ = n, err

	var item2 KV[string, string]
	n, err = sstape.Unroll(bs, &item2)
	_, _ = n, err
	if item != item2 {
		panic("tapeKV works wrongly")
	}
}

func TestKVsi(t *testing.T) {
	sstape := NewTapeKV[string, int](nil, nil)
	item := KV[string, int]{
		Key:   "key",
		Value: 12345678,
	}
	bs := make([]byte, sstape.Sizeof(item))
	n, err := sstape.Roll(item, bs)
	_, _ = n, err

	var item2 KV[string, int]
	n, err = sstape.Unroll(bs, &item2)
	_, _ = n, err
	if item != item2 {
		panic("tapeKV works wrongly")
	}
}

func TestKVis(t *testing.T) {
	sstape := NewTapeKV[int, string](nil, nil)
	item := KV[int, string]{
		Key:   12345678,
		Value: "Value",
	}
	bs := make([]byte, sstape.Sizeof(item))
	n, err := sstape.Roll(item, bs)
	_, _ = n, err

	var item2 KV[int, string]
	n, err = sstape.Unroll(bs, &item2)
	_, _ = n, err
	if item != item2 {
		panic("tapeKV works wrongly")
	}
}

func TestKVst(t *testing.T) {
	sstape := NewTapeKV[string, time.Time](nil, &fastape.TimeTape{})
	item := KV[string, time.Time]{
		Key:   "12345678",
		Value: time.Now(),
	}
	bs := make([]byte, sstape.Sizeof(item))
	n, err := sstape.Roll(item, bs)
	_, _ = n, err

	var item2 KV[string, time.Time]
	n, err = sstape.Unroll(bs, &item2)
	_, _ = n, err
	if item.Value.UTC() != item2.Value.UTC() {
		panic("tapeKV works wrongly")
	}
}
