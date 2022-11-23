package internal

// Argument, Option, Command, Parent
const NumberKinds = 4

const DiffFromValueKind = 200

type FieldKind int

// Field kind
const (
	UnknownKind FieldKind = iota + DiffFromValueKind
	ArgumentKind
	OptionKind
	CommandKind
	ParentKind
)

const (
	UnknownStr  = "Unknown"
	ArgumentStr = "Argument"
	OptionStr   = "Option"
	CommandStr  = "Command"
	ParentStr   = "Parent"
)

func (f FieldKind) KindIndexKey() int {
	return int(f) - DiffFromValueKind
}

func (f FieldKind) String() string {
	switch f {
	case ArgumentKind:
		return ArgumentStr
	case OptionKind:
		return OptionStr
	case CommandKind:
		return CommandStr
	case ParentKind:
		return ParentStr
	default:
		return UnknownStr
	}
}

func FieldKindFromString(s string) FieldKind {
	switch s {
	case ArgumentStr:
		return ArgumentKind
	case OptionStr:
		return OptionKind
	case CommandStr:
		return CommandKind
	case ParentStr:
		return ParentKind
	default:
		return UnknownKind
	}

}
