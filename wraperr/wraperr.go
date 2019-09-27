package wraperr

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

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

func ResultVar(errVar *error, params ...interface{}) {
	if *errVar == nil {
		return
	}
	*errVar = WithStackSkip(1, *errVar, params...)
}

func WithStack(err error, params ...interface{}) error {
	return &withStack{
		err:    err,
		params: params,
		stack:  callStack(1),
	}
}

func WithStackSkip(skip int, err error, params ...interface{}) error {
	return &withStack{
		err:    err,
		params: params,
		stack:  callStack(1 + skip),
	}
}

type withStack struct {
	err    error
	params []interface{}
	stack  []uintptr
}

func (w *withStack) Error() string {
	var b strings.Builder
	fmt.Fprintln(&b, w.err.Error())
	fmt.Println(w.stack)
	frames := runtime.CallersFrames(w.stack)
	frame, ok := frames.Next()
	if !ok {
		panic("no stack frame")
	}
	fmt.Fprintf(&b, "%s(%s)\n\t%s:%d\n", frame.Func.Name(), strings.Join(formatArgs(w.params), ", "), frame.File, frame.Line)
	for ok {
		fmt.Fprintf(&b, "%s\n\t%s:%d\n", frame.Func.Name(), frame.File, frame.Line)
		frame, ok = frames.Next()
	}
	return b.String()
}

func (w *withStack) Unwrap() error {
	return w.err
}

func callStack(skip int) []uintptr {
	c := make([]uintptr, 32)
	n := runtime.Callers(skip+2, c)
	return c[:n]
}
