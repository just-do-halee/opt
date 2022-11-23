package internal

import "fmt"

type BuilderOptions struct {
	Author  string
	Version string
	About   string
}

func (b BuilderOptions) String() string {
	return fmt.Sprint(b.Version, "  ", b.Author, "\n\n", b.About)
}
