package basic

type StringArray struct {
	*BytesArray
}

// func (sa *StringArray) UnsafeGet(index int) *string {
// 	pbs := sa.BytesArray.UnsafeGet(index)
// 	if pbs == nil {
// 		return nil
// 	}
// 	return Ptr(b2s(*pbs))
// }

// unsafe-mode; for later use must be cloned
func (sa *StringArray) Get(index int) string {
	return b2s(sa.UnsafePtr(index))
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

func (sa *StringArray) Seq(yield func(string) bool) {
	for i := range sa.Len() {
		if !yield(b2s(sa.UnsafePtr(i))) {
			break
		}
	}
}
