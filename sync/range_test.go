package sync

import (
	"fmt"
	"testing"
)

func TestMapRangeInt(t *testing.T) {
	m := NewMap[int, int](5)
	for i := 1; i < 20; i++ {
		m.Set(i, -i)
	}
	done := make(chan struct{})
	f := func(id int) {
		m.Range(func(key, value int) bool {
			fmt.Println(id, key, value)
			return true
		})
		done <- struct{}{}
	}
	go f(10)
	go f(100)
	<-done
	<-done
}

func TestMapRangeString(t *testing.T) {
	m := NewMap[string, string](5)
	for i := 1; i < 10; i++ {
		m.Set(fmt.Sprint(i), fmt.Sprint(-i))
	}
	done := make(chan struct{})
	f := func(id int) {
		m.Range(func(key, value string) bool {
			fmt.Println(id, key, value)
			return true
		})
		done <- struct{}{}
	}
	go f(10)
	go f(100)
	<-done
	<-done
}
