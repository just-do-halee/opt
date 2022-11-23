package internal

// FieldOpt is an information of [FieldInfo] struct.
type FieldOpt struct {
	required   bool
	short      bool
	long       bool
	occurrency bool
}

func (o FieldOpt) String() string {
	var sb Stringify
	if o.required {
		sb.WriteString("required,")
	}
	if o.short {
		sb.WriteString("short,")
	}
	if o.long {
		sb.WriteString("long,")
	}
	if o.occurrency {
		sb.WriteString("occurrency,")
	}
	if sb.Len() == 0 {
		return ""
	}
	return sb.String()[:sb.Len()-1]
}

func (o FieldOpt) Required() bool {
	return o.required
}

func (o FieldOpt) hasShortNotLong() bool {
	return o.short && !o.long
}

func (o FieldOpt) hasLongNotShort() bool {
	return o.long && !o.short
}

func (o FieldOpt) hasShortAndLong() bool {
	return o.short && o.long
}

func (o FieldOpt) hasNotShortAndLong() bool {
	return !o.short && !o.long
}
