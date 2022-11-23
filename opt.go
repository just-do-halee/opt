package opt

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/just-do-halee/opt/internal"
	"github.com/just-do-halee/opt/pkg/constraints"
)

//go:inline
func preprocessArgs(args []string) ([]string, error) {
	// preprocessing
	args_ := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch len(arg) {
		case 0:
			continue
		case 1:
			if arg[0] == '-' {
				return nil, errors.New("invalid argument: " + arg)
			}
			fallthrough
		case 2, 3:
			args_ = append(args_, arg)
			continue
		default:
			// -a=b or --a=b
			if arg[0] == '-' {
				argName, argValue, ok := strings.Cut(arg, "=")
				if ok {
					args_ = append(args_, argName, argValue)
				} else {
					args_ = append(args_, arg)
				}
			} else {
				args_ = append(args_, arg)
			}
			continue
		}

	}

	return args_, nil
}

//go:inline
func Args(args ...string) *builder {
	if len(args) == 0 {
		args = os.Args[1:]
	}
	args, err := preprocessArgs(args)
	return &builder{args: args, earlyError: err}
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
