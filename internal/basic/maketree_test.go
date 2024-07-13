package basic

import (
	"fmt"
	"testing"

	"github.com/nazarifard/bigtype/options"
)

func TestMakeTree(t *testing.T) {
	m := NewMap[int, string](20)
	m.Set(1, "1000")
	m.Set(2, "2000")
	m.Set(3, "3000")
	m.Set(4, "4000")
	m.Range(func(key int, value string) bool {
		fmt.Println(key, value)
		return true
	})
}

func TestMakeMap(t *testing.T) {
	var mo options.MapOptions[string, string]
	//m0:=kv.NewTapeKV[string,string](fastape.StringTape{})
	//mo.WithMarshal(fastape.StringTape{})
	m := NewMap[string, string](mo)
	m.Set("1", "1000")
	m.Set("2", "2000")
	m.Set("3", "3000")
	m.Set("4", "4000")
	m.Range(func(key, value string) bool {
		fmt.Println(key, value)
		return true
	})
}
