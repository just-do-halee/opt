package opt

import (
	"errors"
	"fmt"
	"os"

	"github.com/just-do-halee/opt/internal"
	"github.com/just-do-halee/opt/pkg/constraints"
)

//go:inline
func Args(args ...string) *builder {
	if len(args) == 0 {
		args = os.Args[1:]
	}
	return &builder{args: args}
}

func Set[I constraints.InputType, T internal.Setter[I]](field T, value I) {
	field.Set(value)
}

func Unset(field internal.Unsetter) {
	field.Unset()
}

func Validate[I constraints.InputType, T internal.Validater[I]](field T, validator func(I) error) error {
	err := field.Validate(validator)
	if err != nil {
		s := fmt.Sprint("\n  ", err.Error(), "\n\n  ", field.GetUsage(), "\n")
		return errors.New(s)
	}
	return nil
}
