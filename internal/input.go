package internal

import (
	"github.com/just-do-halee/opt/pkg/constraints"
	"github.com/just-do-halee/refl"
)

type Setter[T constraints.InputType] interface {
	Set(value T)
}

type Getter[T constraints.InputType] interface {
	Get() T
}

type Unsetter interface {
	Unset()
}

type Validater[T constraints.InputType] interface {
	Validate(validator func(T) error) error
	GetUsage() string
	SetUsage(usage string)
}

type ValueKindGetter interface {
	ValueKind() refl.Kind
}

type Input[T constraints.InputType] struct {
	value T
	isSet bool
	usage string
}

func (i *Input[T]) GetUsage() string {
	return i.usage
}

func (i *Input[T]) SetUsage(usage string) {
	i.usage = usage
}

func (i *Input[T]) ValueKind() refl.Kind {
	return refl.TypeOf(i.value).Kind()
}

func (i *Input[T]) Set(value T) {
	i.value = value
	i.isSet = true
}

func (i *Input[T]) Unset() {
	var value T
	i.value = value
	i.isSet = false
}

func (i *Input[T]) Get() T {
	return i.value
}

func (i *Input[T]) Validate(validator func(T) error) error {
	return validator(i.value)
}
