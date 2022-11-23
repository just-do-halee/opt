package internal

import (
	"strings"

	"github.com/just-do-halee/refl"
)

func MakeTree(structPtr any, parent any, options BuilderOptions) (tree *Tree, err error) {
	// log.Println("WalkThroughStruct", refl.TypeOf(structPtr), refl.TypeOf(parent))

	if parent != nil && !IsPointer(parent) {
		OptFatal("[OPT] parent must be a pointer")
	}

	// get the detail of the struct
	st := refl.GetStruct(structPtr)
	if st.Size == 0 {
		return
	}

	// run Before() if there is
	if err = RunBefore(structPtr); err != nil {
		return
	}

	// make tree
	// {
	//		Arguments []*FieldInfo
	//		Options   map[string]*FieldInfo
	//		Commands  map[string]*FieldInfo
	//		Required  map[*FieldInfo]bool
	// }
	tree = NewTree(options)

	// make kind indexes to order the tree fields
	kindIndexes := make(map[int]int, NumberKinds)

	// iterate through the fields of the struct
	for i := 0; i < st.Size; i++ {
		field := st.Field(i)

		// check the field unexported
		if field.PkgPath != "" {
			OptFatal(
				"\n[OPT] opt field must be exported\n\n\t",
				field.Parent.Type.Name(),
				" {\n\n\t  ", field.Name, " -> ",
				strings.ToUpper(string(field.Name[0]))+field.Name[1:],
				"\n\n\t}\n\n")
		}

		info, ok := makeFieldInfo(field)
		if !ok {
			// skip the unknown field
			continue
		}

		kind := info.Kind

		if kind == ParentKind {
			if parent == nil {
				OptFatal("[OPT] parent is nil")
			}

			// set the parent field
			field.Value.Addr().Interface().(UnsafeSetter).UnsafeSet(parent)
			continue
		} else {
			// count the proper kinds
			tree.total++

		}

		info.Index = i
		info.KindIndex = kindIndexes[int(kind)]

		// increment the current kind index
		kindIndexes[int(kind)]++

		// add the field info to the tree
		tree.Add(info)
	}

	return
}

const fieldValueName = "value"

func getFieldDetail(field refl.StructField) (kind FieldKind, innerType refl.Type) {
	ut := refl.UnwrapType(field.Type)

	origin, _ := refl.GetTypeName(ut)
	kind = FieldKindFromString(origin)

	if kind == UnknownKind {
		return
	}

	valueType, ok := ut.FieldByName(fieldValueName)
	if !ok {
		kind = UnknownKind
	}
	innerType = refl.UnwrapType(valueType.Type)
	return
}

func makeFieldInfo(field refl.StructField) (info *FieldInfo, ok bool) {
	kind, innerType := getFieldDetail(field)
	if kind == UnknownKind {
		return
	}
	// the field info
	info = &FieldInfo{
		Field:     &field,
		Name:      field.Name,
		Kind:      kind,
		InnerType: innerType,
		Message:   field.Tag.Get(TagMsg), // the tag `msg:"..."`
	}
	if tag := field.Tag.Get(TagOpt); len(tag) > 0 {
		info.Opt = ParseTagOpt(tag)
	}
	ok = true
	return
}
