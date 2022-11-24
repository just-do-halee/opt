package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/goccy/go-reflect"
	"github.com/just-do-halee/refl"
	"github.com/just-do-halee/refl/kind"
)

type UnsafeSetter interface {
	UnsafeSet(value any)
}

//go:inline
func MatchAndSetField(structPtr any, tree *Tree, args *[]string) (nextStruct any, err error) {
	argsLen := len(*args)

	switch {
	// the tree has no any field or no args and no required fields anymore
	case tree.total == 0 || (argsLen == 0 && len(tree.Required) == 0):
		// run After() and Run()
		if err = RunAfter(structPtr); err != nil {
			return
		}
		err = RunRun(structPtr)
		return

	case argsLen == 0:
		// no args, but there are required fields
		// show help
		err = errors.New(tree.ToHelp())
		return

	}

	arguments := tree.Arguments

	// match and set the field
	for i := 0; i < argsLen; i++ {

		// arg = full string name
		arg, argValue, valuedOk := strings.Cut((*args)[i], "=")
		if arg == "" {
			continue
		}

		// argName = pure name (cutted if it's a short name)
		argName := arg
		_, count, e := countShortFlagOccurency(argName)
		if e == nil {
			// if the arg is a short flag, cut the flag name
			argName = arg[:2]
		}

		// option
		if option := tree.Options[argName]; option != nil {

			// -... pure occurrences
			if count > 1 {
				if !option.Opt.occurrency {
					s := fmt.Sprint("\n  invalid occurrences: ", argName, "...\n\n  ", option.Usage, "\n")
					err = errors.New(s)
					return
				}
				if valuedOk {
					s := fmt.Sprint("\n  invalid value in short occurrences: ", argValue, "\ttry ", arg, "\n\n  ", option.Usage, "\n")
					err = errors.New(s)
					return
				}
			}

			countStr := strconv.Itoa(count)

			// single short with the value
			if count == 1 && valuedOk {
				countStr = argValue
			}

			// short
			if count > 0 {
				// e.g. -v..
				// set the value to the field
				err = parsedSet(option, countStr)
				if err != nil {
					return
				}

				// remove the field from the required list
				delete(tree.Required, option)
				continue
			}

			// long with occurrences
			if option.Opt.occurrency && !valuedOk {
				// e.g. --verbose
				// set the value to the field
				err = parsedSet(option, "1")
				if err != nil {
					return
				}

				// remove the field from the required list
				delete(tree.Required, option)
				continue
			}

			if option.IsInnerTypeHelp() && !valuedOk {
				err = errors.New(tree.ToHelp())
				return
			}

			if option.InnerType.Kind() == kind.Bool {

				value := "1"

				if valuedOk {
					_, err = strconv.ParseBool(argValue)
					if err != nil {
						s := fmt.Sprint("\n  invalid value: ", argValue, "\ttry ", arg, "\n\n  ", option.Usage, "\n")
						err = errors.New(s)
						return
					}
					value = argValue
				}

				err = parsedSet(option, value)
				if err != nil {
					return
				}

				delete(tree.Required, option)
				continue
			}

			// normal set
			if valuedOk {
				err = parsedSet(option, argValue)
				if err != nil {
					s := fmt.Sprint("\n  ", err.Error(), "\n\n  ", option.Usage, "\n")
					err = errors.New(s)
					return
				}
			} else {
				if i+1 < len(*args) {
					err = parsedSet(option, (*args)[i+1])
					if err != nil {
						s := fmt.Sprint("\n  ", err.Error(), "\n\n  ", option.Usage, "\n")
						err = errors.New(s)
						return
					}
					// increment i to skip the value
					i++
				} else {
					s := fmt.Sprint("\n  missing value for ", argName, "\n\n  ", option.Usage, "\n")
					err = errors.New(s)
					return
				}
			}

			// remove the field from the required list
			delete(tree.Required, option)
			continue

			// command
		} else if command := tree.Commands[argName]; command != nil {
			if command.IsInnerTypeHelp() {
				err = errors.New(tree.ToHelp())
				return
			}

			if i+1 <= len(*args) {
				*args = (*args)[i+1:]
			} else {
				*args = (*args)[:0]
			}

			nextStruct = reflect.New(command.InnerType).Interface()
			return
			// argument
		} else {

			if len(arguments) == 0 {
				s := fmt.Sprint("\n  unexpected argument \"", argName, "\"\n\n  Use '", tree.Footer(), "' for more information\n")
				err = errors.New(fmt.Sprint(s))
				return
			}

			err = parsedSet(arguments[0], arg)
			if err != nil {
				return
			}

			// remove the field from the required list
			delete(tree.Required, arguments[0])
			arguments = arguments[1:]
		}

	}

	// check if there are any required fields
	if len(tree.Required) > 0 {
		for requiredInfo := range tree.Required {
			s := fmt.Sprint("\n  missing required\t", requiredInfo.Usage, "\n")
			err = errors.New(s)
			return
		}
	}

	// run After() and Run()
	if err = RunAfter(structPtr); err != nil {
		return
	}
	err = RunRun(structPtr)
	return
}

func parsedSet(info *FieldInfo, s string) error {
	inputer := info.Field.Value.Addr().Interface()

	var parsedValue refl.Value
	switch info.InnerType.Kind() {
	case kind.String:
		parsedValue = refl.ValueOf(s)

	case kind.Bool:
		parsed, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		parsedValue = refl.ValueOf(parsed)

	case kind.Int:
		parsed, err := strconv.ParseInt(s, 0, 0)
		if err != nil {
			return err
		}
		parsedValue = refl.ValueOf(int(parsed))

	case kind.Int8:
		parsed, err := strconv.ParseInt(s, 0, 8)
		if err != nil {
			return err
		}
		parsedValue = refl.ValueOf(int8(parsed))

	case kind.Int16:
		parsed, err := strconv.ParseInt(s, 0, 16)
		if err != nil {
			return err
		}
		parsedValue = refl.ValueOf(int16(parsed))

	case kind.Int32:
		parsed, err := strconv.ParseInt(s, 0, 32)
		if err != nil {
			return err
		}
		parsedValue = refl.ValueOf(int32(parsed))

	case kind.Int64:
		parsed, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return err
		}
		parsedValue = refl.ValueOf(int64(parsed))

	case kind.Uint:
		parsed, err := strconv.ParseUint(s, 0, 0)
		if err != nil {
			return err
		}
		parsedValue = refl.ValueOf(uint(parsed))

	case kind.Uint8:
		parsed, err := strconv.ParseUint(s, 0, 8)
		if err != nil {
			return err
		}
		parsedValue = refl.ValueOf(uint8(parsed))

	case kind.Uint16:
		parsed, err := strconv.ParseUint(s, 0, 16)
		if err != nil {
			return err
		}
		parsedValue = refl.ValueOf(uint16(parsed))

	case kind.Uint32:
		parsed, err := strconv.ParseUint(s, 0, 32)
		if err != nil {
			return err
		}
		parsedValue = refl.ValueOf(uint32(parsed))

	case kind.Uint64:
		parsed, err := strconv.ParseUint(s, 0, 64)
		if err != nil {
			return err
		}
		parsedValue = refl.ValueOf(uint64(parsed))

	case kind.Float32:
		parsed, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return err
		}
		parsedValue = refl.ValueOf(float32(parsed))

	case kind.Float64:
		parsed, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		parsedValue = refl.ValueOf(float64(parsed))

	default:
		panic("not implemented")

	}
	parsedValue = parsedValue.Convert(info.InnerType)
	refl.ValueOf(inputer).MethodByName("Set").Call([]refl.Value{parsedValue})
	return nil
}

func countShortFlagOccurency(arg string) (first rune, count int, err error) {
	if arg == "" {
		err = errors.New("empty string")
		return
	}
	if arg[0] != '-' {
		err = errors.New("not a flag")
		return
	}
	if arg[1] == '-' {
		err = errors.New("not a short flag")
		return
	}
	for i := 1; i < len(arg); i++ {
		if arg[i] == arg[1] {
			count++
		} else {
			err = errors.New("not a short occurency flag")
			return
		}
	}
	first = rune(arg[1])
	return
}
