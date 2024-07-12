package basic

// //FileSystem Tree Path implemention over Array
// //Not tested carefully yet

// import (
// 	"github.com/nazarifard/fastape"
// 	"github.com/nazarifard/syncpool"
// )

// type Index = int //[8]byte
// type FsInterface interface {
// 	Set(path [8]byte, data Index)           //write to file
// 	Get(path [8]byte) (data Index, ok bool) //read file data
// }
// type Name = byte
// type Type = bool
// type Data = Index

// const FolderType bool = false
// const FileType bool = true

// type NameIndexTape struct {
// 	NameTape  fastape.UnitTape[Name]
// 	IndexTape fastape.UnitTape[Index]
// }

// func (t *NameIndexTape) Sizeof(n NameIndex) int {
// 	return t.NameTape.Sizeof(n.Name) +
// 		t.IndexTape.Sizeof(n.Index)
// }

// func (t *NameIndexTape) Roll(node NameIndex, bs []byte) (n int, err error) {
// 	m, err := t.NameTape.Roll(node.Name, bs[n:])
// 	if err != nil {
// 		return 0, err
// 	}
// 	n += m

// 	m, err = t.IndexTape.Roll(node.Index, bs[n:])
// 	if err != nil {
// 		return 0, err
// 	}
// 	n += m
// 	return
// }

// func (t *NameIndexTape) Unroll(bs []byte, node *NameIndex) (n int, err error) {
// 	m, err := t.NameTape.Unroll(bs[n:], &node.Name)
// 	if err != nil {
// 		return 0, err
// 	}
// 	n += m

// 	m, err = t.IndexTape.Unroll(bs[n:], &node.Index)
// 	if err != nil {
// 		return 0, err
// 	}
// 	n += m
// 	return
// }

// type NodeTape struct {
// 	TypeTape fastape.UnitTape[Type]
// 	//NameTape     fastape.UnitTape[Name]
// 	ChildrenTape fastape.SliceTape[NameIndex, fastape.UnitTape[NameIndex]]
// }

// func (t *NodeTape) Sizeof(n Node) int {
// 	return t.TypeTape.Sizeof(n.Type) +
// 		t.ChildrenTape.Sizeof(n.Children)
// }

// func (t *NodeTape) Roll(node Node, bs []byte) (n int, err error) {
// 	m, err := t.TypeTape.Roll(node.Type, bs)
// 	if err != nil {
// 		return 0, err
// 	}
// 	n += m

// 	m, err = t.ChildrenTape.Roll(node.Children, bs[n:])
// 	if err != nil {
// 		return 0, err
// 	}
// 	n += m
// 	return
// }

// func (t *NodeTape) Unroll(bs []byte, node *Node) (n int, err error) {
// 	m, err := t.TypeTape.Unroll(bs[n:], &node.Type)
// 	if err != nil {
// 		return 0, err
// 	}
// 	n += m

// 	m, err = t.ChildrenTape.Unroll(bs[n:], &node.Children)
// 	if err != nil {
// 		return 0, err
// 	}
// 	n += m
// 	return
// }

// type nodeTap struct {
// 	*NodeTape
// 	syncpool.BufferPool
// }

// var NodeTap = &nodeTap{
// 	NodeTape:   &NodeTape{},
// 	BufferPool: syncpool.NewBufferPool(),
// }

// func (nt *nodeTap) Encode(n Node) (syncpool.Buffer, error) {
// 	size := nt.NodeTape.Sizeof(n)
// 	buf := nt.BufferPool.Get(size)
// 	_, err := nt.NodeTape.Roll(n, buf.Bytes())
// 	if err != nil {
// 		buf.Free()
// 	}
// 	return buf, err
// }

// func (nt *nodeTap) Decode(bs []byte) (node Node, n int, err error) {
// 	n, err = nt.NodeTape.Unroll(bs, &node)
// 	return
// }

// type NameIndex struct {
// 	Name
// 	Index
// }

// type Node struct {
// 	Type
// 	Children []NameIndex //[]*Node
// }

// type File Node

// func (f *File) Write(data Data) {
// 	switch len(f.Children) {
// 	case 0:
// 		f.Children = append(f.Children, NameIndex{Index: data})
// 	case 1:
// 		f.Children[0] = NameIndex{Index: data}
// 	default:
// 		panic("seemingly file content is wrong")
// 	}
// }

// func (f *File) Read() (data Data, ok bool) {
// 	switch len(f.Children) {
// 	case 0:
// 		return //invalid
// 	case 1:
// 		return f.Children[0].Index, true
// 	default:
// 		return //invalid
// 	}
// }

// // func NewHard() *Hard {
// // 	var h Hard
// // 	for i := range len(h.Disk) {
// // 		h.Disk[i].Array = NewArray[Node](0, NodeTap, true)
// // 		h.Disk[i].Array.Set(0, Node{})
// // 	}
// // 	return &h
// // }
// // type Hard struct {
// // 	Disk [256 * 256]Partition
// // }
// // func (h *Hard) Set(id uint64, data Index) {
// // 	i := uint16(id) & 0xFFFF
// // 	h.Disk[i].Set(id<<16, data)
// // }
// // func (h *Hard) Get(id uint64) (data Data, ok bool) {
// // 	i := uint16(id) & 0xFFFF
// // 	return h.Disk[i].Get(id << 16)
// // }

// type Partition struct {
// 	Array[Node]
// }

// func (d *Partition) Mkdir(path Index, name Name) (i Index, ok bool) {
// 	if path > d.Array.Len() {
// 		return
// 	}
// 	parent := d.Array.Get(path)
// 	child := Node{
// 		Type:     FolderType,
// 		Children: nil,
// 	}
// 	index := d.Array.Len()
// 	d.Array.Set(index, child)
// 	parent.Children = append(parent.Children, NameIndex{Index: index})
// 	d.Array.Set(path, parent) //update parent
// 	return index, true
// }

// func (d *Partition) Touch(path Index, name Name, data Data) Index {
// 	folder := d.Array.Get(path)
// 	file := Node{
// 		Type:     FileType,
// 		Children: []NameIndex{{Name: name, Index: data}},
// 	}
// 	index := d.Array.Len()
// 	d.Array.Set(index, file)
// 	folder.Children = append(folder.Children, NameIndex{Name: name, Index: index})
// 	d.Array.Set(path, folder) //update parent
// 	return index
// }

// func (d *Partition) Len() int {
// 	return d.Array.Len()
// }

// func (d *Partition) Set(id uint64, data Index) {
// 	var root Index = 0 //0 Invalid
// 	var i int
// 	var found bool = false
// 	for i = 0; i < 7; i++ { //sizeof uint64
// 		name := byte(id & 0xFF)
// 		folders := d.Ls(root)
// 		found = false
// 		for _, f := range folders { //index is second parameter not first! here index means node pointer
// 			if f.Name == name {
// 				id >>= 8
// 				root = f.Index
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			break
// 		}
// 	}
// 	var ok bool
// 	for j := i; j < 7; j++ {
// 		name := byte(id & 0xFF)
// 		id >>= 8
// 		root, ok = d.Mkdir(root, name)
// 		if !ok {
// 			return
// 		}
// 	}
// 	d.Touch(root, byte(id), data)
// }

// func (d *Partition) Get(id uint64) (data Data, ok bool) {
// 	var root Index = 0
// ROOT:
// 	for range 6 {
// 		name := byte(id & 0xFF)
// 		folders := d.Ls(root)
// 		for _, f := range folders {
// 			if f.Name == name {
// 				root = f.Index
// 				continue ROOT
// 			}
// 		}
// 		return //not found
// 	}
// 	f := d.Array.Get(root)
// 	return (*File)(&f).Read()
// }

// func (d *Partition) Ls(root Index) []NameIndex {
// 	if d.Array.Len() < root {
// 		return nil
// 	}
// 	return d.Array.Get(root).Children
// }
