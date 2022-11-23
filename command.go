package opt

import (
	"github.com/just-do-halee/opt/internal"
)

type Command[T any] struct {
	value *T
	ok    bool
	usage string
}

func (c *Command[T]) Ok() bool {
	return c.ok
}

func (c *Command[T]) GetUsage() string {
	return c.usage
}

func (c *Command[T]) SetUsage(usage string) {
	c.usage = usage
}

func (c *Command[T]) GetPtr() any {
	return c.value
}

func (c *Command[T]) UnsafeSet(value any) {
	c.value = value.(*T)
}

func (c *Command[T]) RunCommand() error {
	return internal.RunRun(c.value)
}
