package bigtype

import (
	"errors"
	"unsafe"
)

var errKeyNotFound = errors.New("key not found")

type Cache struct {
	memory Map[string, []byte]
}

func NewCache() *Cache {
	return &Cache{
		memory: NewMap[string, []byte](),
	}
}

func (c *Cache) Get(key []byte) (value []byte, err error) {
	kStr := unsafe.String(&key[0], len(key))
	value, ok := c.memory.Get(kStr)
	if !ok {
		return nil, errKeyNotFound
	}
	return
}

func (c *Cache) Set(key, value []byte) error {
	kStr := unsafe.String(&key[0], len(key))
	c.memory.Set(kStr, value)
	return nil
}
