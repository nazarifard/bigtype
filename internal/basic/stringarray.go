package basic

import "strings"

type StringArray struct {
	*BytesArray
}

func (sa *StringArray) UnsafePtr(index int) *string {
	pbs := sa.BytesArray.UnsafePtr(index)
	if pbs == nil {
		return nil
	}
	return Ptr(b2s(*pbs))
}

func (sa *StringArray) Get(index int) string {
	pStr := sa.UnsafePtr(index)
	if pStr == nil {
		return ""
	}
	return strings.Clone(*pStr)
}

func (sa *StringArray) Set(index int, value string) {
	sa.BytesArray.Set(index, s2b(value))
}

func (ba *StringArray) Update(index int, fn func(string) string) {
	ba.BytesArray.Update(index, func(old []byte) []byte {
		return s2b(fn(b2s(old)))
	})
}

func NewStringArray(size int, expandable bool) *StringArray {
	return &StringArray{
		BytesArray: NewBytesArray(size, expandable),
	}
}
