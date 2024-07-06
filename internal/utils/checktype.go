package utils

import "reflect"

func IsFixedType(v interface{}) bool {
	return isFixedType(reflect.TypeOf(v))
}

func isFixedType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Slice,
		reflect.Map,
		reflect.Chan,
		reflect.String,
		reflect.Func,
		reflect.Interface,
		reflect.Invalid,
		reflect.Pointer,
		reflect.UnsafePointer:
		return false

	case reflect.Array:
		return isFixedType(t.Elem())

	case reflect.Struct:
		len := t.NumField()
		for i := 0; i < len; i++ {
			if !isFixedType(t.Field(i).Type) {
				return false
			}
		}
	}
	return true
}

func IsString(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.String
}

func IsBytes(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice &&
		reflect.TypeOf(v).Elem().Kind() == reflect.Uint8
}
