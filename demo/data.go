package main

import (
	"fmt"
	"time"

	"github.com/nazarifard/bigtype"
	"github.com/sanity-io/litter"
)

type Person struct {
	IsSingle     bool
	Age          int
	Name         string
	Addresses    []string
	BirthDate    *time.Time
	CustomFields map[string]string
}

func demoArray() {
	type Sample struct {
		Num      uint32
		spId     [12]byte
		simId    [12]byte
		spId2    [12]byte
		mobileId [12]byte
		//comment [200]byte
		isValid bool
	}
	id := [12]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2'}
	value := Sample{1_234_567, id, id, id, id, false}

	size := 2_000_000
	index := 1_234_567

	fixedArr := bigtype.NewArray[Sample](int(size)) //Sample is a fixed sized data type
	fixedArr.Set(index, value)
	v := fixedArr.Get(index)
	if v != value {
		fmt.Println("Error: BigFixedArray doesn't work properly.")
	}

	arr := bigtype.NewArray[string](int(size)) //string is not a fixed sized data type
	str := "abcd_abcd_abcd"
	arr.Set(index, str)
	s := arr.Get(index)
	if s != str {
		fmt.Println("Error: BigArray doesn't work properly.")
	}

	fmt.Println("BigArray works well.")
}

func demoMap() {
	//var m model.Map[string, string]
	const maxSize = 1_000

	m := bigtype.NewMap[string, string](maxSize)
	m.Set("aaa", "111")
	v, ok := m.Get("aaa")
	if !ok || v != "111" {
		fmt.Println("Error: NewBigmap doesnt work fine")
	}

	//type MString = fastape.StringTape
	//m2 := bigtype.NewBigmapWithM[string, Person, MString, PersonTape](maxSize)
	m2 := bigtype.NewMap[string, Person](maxSize)

	birthDate, _ := time.Parse("2006-01-02", "1980-03-22")
	trueman := Person{
		IsSingle:     true,
		Age:          32,
		Name:         "Babi",
		Addresses:    []string{"Home", "Office1", "Office2"},
		BirthDate:    &birthDate,
		CustomFields: map[string]string{"one": "1", "two": "2"},
	}
	m2.Set("trueman", trueman)
	man, ok := m2.Get("trueman")

	if !ok {
		fmt.Println("Error: NewBigmapWithM doesnt work fine")
	}

	fmt.Println("\ntrueman: ")
	litter.Dump(trueman)

	fmt.Println("\nman: ")
	litter.Dump(man)
}

func test_demo() {
	const SIZE = 1_000_000
	bigmap := bigtype.NewMap[string, uint64]()

	start := time.Now()
	for i := range uint64(SIZE) {
		bigmap.Set(fmt.Sprint(i), 10*i)
	}
	fmt.Printf("\n insert time: %v", time.Since(start))

	start = time.Now()
	for i := range uint64(SIZE - 100) {
		bigmap.Set(fmt.Sprint(i), i+1)
		j, ok := bigmap.Get(fmt.Sprint(i))
		if !ok || j != i+1 {
			fmt.Printf("i!=-j %v!=%v, ok: %v", i, j, ok)
			//os.Exit(1)
		}
	}
	fmt.Printf("\n update time: %v", time.Since(start))

	start = time.Now()
	for i := uint64(0); i < uint64(2); i++ {
		j, ok := bigmap.Get(fmt.Sprint(i))
		fmt.Println(i, j, ok)
		if !ok || j != i+1 {
			fmt.Printf("i!=-j %v!=%v, ok: %v", i, j, ok)
			//os.Exit(1)
		}
	}
	fmt.Printf("\n search time: %v", time.Since(start))

	fmt.Println()

}
