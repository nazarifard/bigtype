package main

import (
	"fmt"

	"github.com/nazarifard/bigtype"
	"github.com/nazarifard/bigtype/log"
	"github.com/sanity-io/litter"
)

func init() {
	log.SetLogger(log.ZeroLog{})
}

func main() {
	litter.Config.HidePrivateFields = false
	demoMap()
	if false {
		demoArray()
		demoMap()
		test_demo()
	}

	m := bigtype.NewMap[int, string](10_000_000)
	m.Set(123, "abc")
	v, ok := m.Get(123)
	fmt.Printf("BigMap key: 123, value: %v, ok: %v\n", v, ok)

	bitArr := bigtype.NewArray[bool](800_000_000) //realSize=size/8=100M
	const truebit = true
	const index = 123_456_789
	bitArr.Set(index, truebit)
	bit := bitArr.Get(index)
	fmt.Printf("BitArray truebit: %v, bit: %v, index: %v\n", truebit, bit, index)
}
