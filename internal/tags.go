package internal

import (
	"strings"
)

// Tags is a list of tags used in this package.
const (
	TagMsg    = "msg"
	TagOpt    = "opt"
	TagRename = "rename"
)

func ParseTagOpt(tagOpt string) (options FieldOpt) {
	tag := strings.Split(tagOpt, ",")
	for i := range tag {
		tag[i] = strings.TrimSpace(tag[i])
		switch tag[i] {
		// case tags
		case "s", "short":
			options.short = true
		case "l", "long":
			options.long = true
		case "o", "occurrency":
			options.occurrency = true
		}
	}
	return
}
