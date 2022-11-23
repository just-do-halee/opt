package opt

import (
	"github.com/just-do-halee/opt/internal"
	"github.com/just-do-halee/opt/pkg/constraints"
)

type Option[T constraints.InputType] struct {
	internal.Input[T]
}
