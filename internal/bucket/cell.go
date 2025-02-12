package bucket

import (
	"unsafe"
)

type CellStatus = bool

const Dead CellStatus = false
const Live CellStatus = true

var _ = func() {
	var h Head
	_ = map[bool]bool{
		false:                       false,
		unsafe.Sizeof(h.Index) >= 5: true,
	}
}

// TODO Len should be [3]byte, instead index should be uint32
type Head struct {
	Len    uint16     // =sizeof(Header)+len(body)+sizeof(Footer)
	Status CellStatus //64bit
	index  [5]byte    //linked back to main array; 40bit
}

func (h *Head) Index() int {
	n := int(*(*uint32)(unsafe.Pointer(&h.index[0])))
	n |= int(h.index[4]) << 32
	return n
}

func (h *Head) SetIndex(n int) {
	*(*uint32)(unsafe.Pointer(&h.index[0])) = uint32(n)
	h.index[4] = byte(n >> 32)
}

// type Tail struct {
// 	Tlen uint16 //always must be equal with Header.Len
// }

type CellPtr struct {
	headPtr uintptr //*Head
	bodyPtr uintptr //*byte
}

//	type Cell1 struct {
//		Head uintptr //*Head
//		Body uintptr //*byte
//	}
// type Cell2 struct {
// 	CellPtr
// 	*Tail
// }

// func (c Cell2) Sizeof() uint16 {
// 	x := uint16(unsafe.Sizeof(*c.Head))
// 	y := c.Head.Len //uint16(len(c.Body))
// 	z := uint16(unsafe.Sizeof(*c.Tail))
// 	return x + y + z
// }

// func (c *Cell2) MergeNext(next Cell2) {
// 	c.CellPtr.MergeNext(next.CellPtr)
// 	c.Tail = next.Tail
// 	c.Tail.Tlen = c.headPtr.Len
// 	//dont touch other fields
// }

func (c *CellPtr) MergeNext(next CellPtr) {
	Header(c.headPtr).Len += Header(next.headPtr).Len
	//dont touch other fields
}

func Header(h uintptr) *Head {
	return (*Head)(unsafe.Pointer(h))
}

func Body(b uintptr) *byte {
	return (*byte)(unsafe.Pointer(b))
}
