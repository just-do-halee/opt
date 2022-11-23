package opt

import (
	"log"

	"github.com/just-do-halee/refl"
)

type Parent[T any] struct {
	value *T
}

func (p *Parent[T]) Get() *T {
	return p.value
}

func (p *Parent[T]) UnsafeSet(value any) {
	if t, ok := value.(*T); ok {
		p.value = t
		return
	}
	log.Fatal("[OPT] invalid type: ", refl.TypeOf(value))
}
