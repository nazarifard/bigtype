package bucket

import "github.com/nazarifard/bigtype/internal/addr"

type Bucket interface {
	Write(index int, data []byte) (offset uint16)
	Read(offset uint16) (data []byte)
	IsRequiredToClean() bool
	Delete(offset uint16)
	Defrag(ad *addr.AddressTable)
	Deadspace() int
	InfoSize() int
	Id() int
	//Bytes() []byte
}
