package internal

type Beforer interface {
	Before() error
}

type Afterer interface {
	After() error
}

type Runner interface {
	Run() error
}

func RunBefore(i any) error {
	if optSt, ok := i.(Beforer); ok {
		return optSt.Before()
	}
	return nil
}

func RunAfter(i any) error {
	if optSt, ok := i.(Afterer); ok {
		return optSt.After()
	}
	return nil
}

func RunRun(i any) error {
	if optSt, ok := i.(Runner); ok {
		return optSt.Run()
	}
	return nil
}
