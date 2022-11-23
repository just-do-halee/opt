package internal

func Parse(structPtr any, args []string, options BuilderOptions) (tree *Tree, err error) {
	var parent any
	var nextStruct any

	for {

		tree, err = MakeTree(structPtr, parent, options) // run Before() and ...
		if err != nil {
			return
		}

		nextStruct, err = MatchAndSetField(structPtr, tree, &args)
		if err != nil || nextStruct == nil {
			return
		}

		parent = structPtr
		structPtr = nextStruct

	}

}
