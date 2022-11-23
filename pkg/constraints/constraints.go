package constraints

import "github.com/goccy/go-reflect"

type Float interface {
	~float32 | ~float64
}

func IsFloat(t reflect.Type) bool {
	kind := t.Kind()
	return kind == reflect.Float32 || kind == reflect.Float64
}

func IsStrFloat(typeName string) bool {
	return typeName == "float32" || typeName == "float64"
}

type SliceFloat interface {
	~[]float32 | ~[]float64
}

func IsSliceFloat(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && IsFloat(t.Elem())
}

func IsStrSliceFloat(typeName string) bool {
	return typeName == "[]float32" || typeName == "[]float64"
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func IsSigned(t reflect.Type) bool {
	kind := t.Kind()
	return kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64
}

func IsStrSigned(typeName string) bool {
	return typeName == "int" || typeName == "int8" || typeName == "int16" || typeName == "int32" || typeName == "int64"
}

type SliceSigned interface {
	~[]int | ~[]int8 | ~[]int16 | ~[]int32 | ~[]int64
}

func IsSliceSigned(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && IsSigned(t.Elem())
}

func IsStrSliceSigned(typeName string) bool {
	return typeName == "[]int" || typeName == "[]int8" || typeName == "[]int16" || typeName == "[]int32" || typeName == "[]int64"
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func IsUnsigned(t reflect.Type) bool {
	kind := t.Kind()
	return kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 || kind == reflect.Uintptr
}

func IsStrUnsigned(typeName string) bool {
	return typeName == "uint" || typeName == "uint8" || typeName == "uint16" || typeName == "uint32" || typeName == "uint64" || typeName == "uintptr"
}

type SliceUnsigned interface {
	~[]uint | ~[]uint8 | ~[]uint16 | ~[]uint32 | ~[]uint64 | ~[]uintptr
}

func IsSliceUnsigned(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && IsUnsigned(t.Elem())
}

func IsStrSliceUnsigned(typeName string) bool {
	return typeName == "[]uint" || typeName == "[]uint8" || typeName == "[]uint16" || typeName == "[]uint32" || typeName == "[]uint64" || typeName == "[]uintptr"
}

type Integer interface {
	Signed | Unsigned
}

func IsInteger(t reflect.Type) bool {
	return IsSigned(t) || IsUnsigned(t)
}

func IsStrInteger(typeName string) bool {
	return IsStrSigned(typeName) || IsStrUnsigned(typeName)
}

type SliceInteger interface {
	SliceSigned | SliceUnsigned
}

func IsSliceInteger(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && IsInteger(t.Elem())
}

func IsStrSliceInteger(typeName string) bool {
	return IsStrSliceSigned(typeName) || IsStrSliceUnsigned(typeName)
}

type String = string

func IsString(t reflect.Type) bool {
	return t.Kind() == reflect.String
}

func IsStrString(typeName string) bool {
	return typeName == "string"
}

type SliceString = []String

func IsSliceString(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && IsString(t.Elem())
}

func IsStrSliceString(typeName string) bool {
	return typeName == "[]string"
}

type Bool = bool

func IsBool(t reflect.Type) bool {
	return t.Kind() == reflect.Bool
}

func IsStrBool(typeName string) bool {
	return typeName == "bool"
}

type SliceBool = []Bool

func IsSliceBool(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && IsBool(t.Elem())
}

func IsStrSliceBool(typeName string) bool {
	return typeName == "[]bool"
}

type InputType interface {
	Float | Integer | ~String | ~Bool |

		SliceFloat | SliceInteger | SliceString | SliceBool
}

func IsInputType(t reflect.Type) bool {
	return IsFloat(t) || IsInteger(t) || IsString(t) || IsBool(t) || IsSliceFloat(t) || IsSliceInteger(t) || IsSliceString(t) || IsSliceBool(t)
}

func IsStrInputType(typeName string) bool {
	return IsStrFloat(typeName) || IsStrInteger(typeName) || IsStrString(typeName) || IsStrBool(typeName) || IsStrSliceFloat(typeName) || IsStrSliceInteger(typeName) || IsStrSliceString(typeName) || IsStrSliceBool(typeName)
}
