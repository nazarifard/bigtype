package basic

import (
	"fmt"
	"testing"
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
	m := NewMap[string, string](20)
	m.Set("1", "1000")
	m.Set("2", "2000")
	m.Set("3", "3000")
	m.Set("4", "4000")
	m.Range(func(key, value string) bool {
		fmt.Println(key, value)
		return true
	})
}
