package unsafereflect

import (
	"reflect"
	_ "unsafe"
)

type abiType pointerOnly

//go:linkname TypesByString reflect.typesByString
func TypesByString(s string) []*abiType

//go:linkname ToType reflect.toType
func ToType(t *abiType) reflect.Type
