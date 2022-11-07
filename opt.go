package opt

import "os"

func Args(args ...string) Builder {
	if len(args) == 0 {
		return Builder{os.Args[1:]}
	}
	return Builder{args}
}

type Builder struct {
	args []string
}

func (a Builder) Build(opt any) error {
	return nil
}
