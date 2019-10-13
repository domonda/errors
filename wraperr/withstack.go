package wraperr

import (
	"errors"
	"fmt"
	"runtime"
)

func WithStack(err error) error {
	return WithStackSkip(1, err)
}

func WithStackSkip(skip int, err error) error {
	return &withStack{
		err:   err,
		stack: callStack(1 + skip),
	}
}

func New(text string) error {
	return WithStackSkip(1, errors.New(text))
}

func Errorf(format string, a ...interface{}) error {
	return WithStackSkip(1, fmt.Errorf(format, a...))
}

type callStackProvider interface {
	Unwrap() error
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
