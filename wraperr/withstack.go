package wraperr

import (
	"errors"
	"fmt"
	"runtime"
)

func WithStack(err error) error {
	return &withStack{
		err:   err,
		stack: callStack(1),
	}
}

func WithStackSkip(skip int, err error) error {
	return &withStack{
		err:   err,
		stack: callStack(1 + skip),
	}
}

func New(text string) error {
	return &withStack{
		err:   errors.New(text),
		stack: callStack(1),
	}
}

func Errorf(format string, a ...interface{}) error {
	return &withStack{
		err:   fmt.Errorf(format, a...),
		stack: callStack(1),
	}
}

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

func callStack(skip int) []uintptr {
	c := make([]uintptr, 32)
	n := runtime.Callers(skip+2, c)
	return c[:n]
}
