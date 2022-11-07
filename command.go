package opt

// Interface
type Runner interface {
	Run() error
}

type Command[T Runner] struct {
	Opt T
	ok  bool
}

// RunCommand runs the command with the given Opt type.
// The Opt type must implement the Runner interface.
func (c Command[T]) RunCommand() error {
	return c.Opt.Run()
}
