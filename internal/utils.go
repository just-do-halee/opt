package internal

import (
	"strings"

	"github.com/goccy/go-reflect"
	"github.com/just-do-halee/refl"
	"github.com/just-do-halee/refl/kind"
)

type Stringify struct {
	strings.Builder
}

func (sb *Stringify) WriteStrings(s ...string) (int, error) {
	var total int
	for _, str := range s {
		num, err := sb.WriteString(str)
		total += num
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func ToAny(i any) any {
	return i
}

func SetValueInStruct[T any](vStruct reflect.Value, key string, value T) bool {
	field := vStruct.FieldByName(key)
	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(value))
		return true
	}
	return false
}

//go:inline
func IsPointer(a any) bool {
	return refl.TypeOf(a).Kind() == kind.Ptr
}
