package bucket

import "github.com/nazarifard/bigtype/internal/addr"

type Bucket interface {
	//Write(index int, data []byte) (offset uint16)
	Request(index int, Len int) (space []byte, offset uint16)
	Get(offset uint16) (data []byte)
	IsRequiredToClean() bool
	Delete(offset uint16)
	Defrag(ad *addr.AddressTable)
	Deadspace() int
	InfoSize() int
	Id() int
	//Bytes() []byte
}
