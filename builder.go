package opt

import (
	"github.com/just-do-halee/opt/internal"
)

type builder struct {
	options internal.BuilderOptions
	args    []string
}

func (b *builder) Author(author string) *builder {
	b.options.Author = author
	return b
}

func (b *builder) Version(version string) *builder {
	b.options.Version = version
	return b
}

func (b *builder) About(about string) *builder {
	b.options.About = about
	return b
}

func (b *builder) Build(opt any) error {
	_, err := internal.Parse(opt, b.args, b.options)
	return err
}

type Tree = *internal.Tree

func (b *builder) BuildForDebugging(opt any) (error, Tree) {
	tree, err := internal.Parse(opt, b.args, b.options)
	return err, tree
}
