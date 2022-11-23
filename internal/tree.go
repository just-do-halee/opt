package internal

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/goccy/go-reflect"
	"github.com/just-do-halee/refl"
)

// Tree is for parsing and matching arguments.
// It also helps to generate help message.
type Tree struct {
	builderOptions BuilderOptions
	// Total number of atomic arguments, options, and commands
	total     int
	Arguments []*FieldInfo
	Options   map[string]*FieldInfo
	Commands  map[string]*FieldInfo
	Required  map[*FieldInfo]bool
	footer    string
}

func NewTree(options BuilderOptions) *Tree {
	return &Tree{
		builderOptions: options,
		Arguments:      []*FieldInfo{},
		Options:        make(map[string]*FieldInfo),
		Commands:       make(map[string]*FieldInfo),
		Required:       make(map[*FieldInfo]bool),
	}
}

func (tree *Tree) BuilderOptions() BuilderOptions {
	return tree.builderOptions
}

func (tree *Tree) Total() int {
	return tree.total
}

func (tree *Tree) Footer() string {
	length := len(tree.footer)
	if length > 2 {
		return tree.footer[:length-2]
	}
	return tree.footer
}

func (tree *Tree) Add(info *FieldInfo) {
	// validate the field info according to the kind and tags
	if err := info.Validate(); err != nil {
		OptFatal(err)
	}

	// argName that will be used in the usage and parsing time
	argName := info.Name

	// renaming if there is a tag `rename:"..."`
	if rename := info.Field.Tag.Get(TagRename); len(rename) > 0 {
		argName = rename
	}

	// add the field info to the tree
	switch info.Kind {
	case ArgumentKind:
		argName = strings.ToUpper(argName)

		// format the usage
		info.Usage = fmt.Sprintf(

			"[%s %s]\t%s",

			argName,
			strings.ToUpper(info.InnerType.Name()),
			info.Message)

		// add to the arguments
		tree.Arguments = append(tree.Arguments, info)

	case OptionKind:
		var sb Stringify

		if info.Opt.hasShortNotLong() {

			// short name

			argName = getShortArgName(argName)

			// for formatting the usage
			sb.WriteStrings(argName, " ")

			// add to the options with the short name
			tree.Options[argName] = info

		} else if info.Opt.hasLongNotShort() {

			// long name

			argName = getLongArgName(argName)

			// for formatting the usage
			sb.WriteStrings("    ", argName, " ")

			// add to the options with the long name
			tree.Options[argName] = info

		} else {

			// both short and long name

			shortName := getShortArgName(argName)
			longName := getLongArgName(argName)

			// for formatting the usage
			sb.WriteStrings(shortName, ", ", longName, " ")

			// add to the options with the both name
			tree.Options[shortName] = info
			tree.Options[longName] = info

			argName = shortName + ", " + longName
		}

		// if the inner value type is not a boolean type and the opt tag hasn't an occurrency
		if !info.Opt.occurrency && info.InnerType.Kind() != reflect.Bool {

			// add the inner type name to the usage
			sb.WriteStrings(strings.ToUpper(info.InnerType.Name()), " ")

		}

		// add the message to the usage
		sb.WriteStrings("\t", info.Message)

		// set the usage
		info.Usage = sb.String()

	case CommandKind:
		argName = strings.ToLower(argName)

		// format the usage
		info.Usage = fmt.Sprintf(

			"%s\t%s",

			argName,
			info.Message)

		// add to the commands
		tree.Commands[argName] = info

	}

	// set the usage to the field value
	info.Field.Value.Addr().MethodByName("SetUsage").Call([]refl.Value{refl.ValueOf(info.Usage)})

	if info.IsInnerTypeHelp() {
		tree.footer += argName + ", "
	}

	// add to the required map if it is required
	// and add to the default value message at the end of the usage
	switch info.Kind {
	case ArgumentKind, OptionKind:
		if !info.IsInnerTypeHelp() {
			if !info.hasDefaultValue() {
				tree.Required[info] = true
			} else {
				info.Usage += fmt.Sprint(" (default: '", info.InnerGet(), "')")
			}
		}
	}
}

func (tree *Tree) ToHelp() string {
	var sb Stringify

	used := make(map[*FieldInfo]bool)

	// for ordering
	ordereds := [3]batchInfo{
		{
			"Commands",
			make([]*FieldInfo, len(tree.Commands)),
		},
		{
			"Arguments",
			make([]*FieldInfo, len(tree.Arguments)),
		},
		{
			"Options",
			make([]*FieldInfo, len(tree.Options)),
		},
	}

	for _, cmd := range tree.Commands {
		if used[cmd] {
			continue
		}
		ordereds[0].append(cmd)
		used[cmd] = true
	}
	for _, arg := range tree.Arguments {
		if used[arg] {
			continue
		}
		ordereds[1].append(arg)
		used[arg] = true
	}

	for _, flag := range tree.Options {
		if used[flag] {
			continue
		}
		ordereds[2].append(flag)
		used[flag] = true
	}

	writer := tabwriter.NewWriter(&sb, 0, 8, 1, '\t', tabwriter.AlignRight)

	fmt.Fprint(writer, "\n", tree.builderOptions, "\n")

	// put
	for _, ordered := range ordereds {
		if len(ordered.Slice) > 0 {
			fmt.Fprintf(writer, "\n%s:\n", ordered.Name)
		}
		for _, info := range ordered.Slice {
			if info == nil {
				continue
			}
			fmt.Fprintf(writer, "\t%s\n", info.Usage)
		}
	}

	fmt.Fprintln(writer)

	writer.Flush()

	return sb.String()
}

type batchInfo struct {
	Name  string
	Slice []*FieldInfo
}

func (b *batchInfo) append(info *FieldInfo) {
	b.Slice[info.KindIndex] = info
}

func getShortArgName(name string) (argName string) {
	for _, c := range name {
		argName = "-" + strings.ToLower(string(c))
		break
	}
	return
}

func getLongArgName(name string) string {
	return "--" + strings.ToLower(name)
}
