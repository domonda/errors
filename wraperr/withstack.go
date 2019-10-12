package wraperr

type callStackProvider interface {
	CallStack() []uintptr
}

type withStack struct {
	err   error
	stack []uintptr
}

func (w *withStack) Error() string {
	return formatError(w)
}

func (w *withStack) Unwrap() error {
	return w.err
}

func (w *withStack) CallStack() []uintptr {
	return w.stack
}
