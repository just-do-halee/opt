package internal

import (
	"errors"

	"github.com/just-do-halee/opt/pkg/constraints"
	"github.com/just-do-halee/refl"
)

// FieldInfo holds information about a field in a opter struct.
type FieldInfo struct {
	Field     *refl.StructField
	Index     int
	Name      string
	Kind      FieldKind
	KindIndex int
	InnerType refl.Type
	Message   string
	Opt       FieldOpt
	Usage     string
}

func (f *FieldInfo) InnerGet() any {
	if f.Field != nil && (f.Kind == ArgumentKind || f.Kind == OptionKind) {
		return f.Field.Value.Addr().MethodByName("Get").Call([]refl.Value{})[0].Interface()
	}
	return nil
}

func (f *FieldInfo) hasDefaultValue() bool {
	if f.Field != nil && (f.Kind == ArgumentKind || f.Kind == OptionKind) {
		return f.Field.Value.FieldByName("isSet").Bool()
	}
	return false
}
func (f *FieldInfo) IsInnerTypeHelp() bool {
	return f.InnerType.Name() == "Help"
}

func (f *FieldInfo) Validate() (err error) {
	if f.Kind == UnknownKind {
		err = errors.New("unknown field kind")
		return
	}

	makeError := func(msg ...string) error {
		var sb Stringify
		sb.WriteStrings(f.Name, ": ", f.InnerType.Name(), "Type ")
		sb.WriteStrings(msg...)
		return errors.New(sb.String())
	}
	cannotHave := func(o FieldOpt, name string) error {
		oStr := o.String()
		if len(oStr) == 0 {
			return nil
		}
		return makeError(name, " cannot have [", oStr, "]")
	}

	// copy
	o := f.Opt

	switch f.Kind {
	case ArgumentKind:
		// the argument can have a required tag
		if o.required {
			o.required = false
		}
		err = cannotHave(o, "Argument")

	case OptionKind:
		// the option can have
		// required, short, long, occurency
		if o.required {
			o.required = false
		}
		if o.short {
			o.short = false
		}
		if o.long {
			o.long = false
		}
		if o.occurrency {
			if !constraints.IsInteger(f.InnerType) {
				return makeError("occurency option can have only integer type")
			}
			o.occurrency = false
		}
		err = cannotHave(o, "Option")
	case CommandKind:
		err = cannotHave(o, "Command")
	}
	return
}
