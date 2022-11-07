package opt

type Help[T any] struct {
	opt     T
	message string
}

func (h Help[T]) String() string {
	return h.message
}

func (h *Help[T]) Parse() error {
	// a := new(T)
	// h.message =
	return nil
}
