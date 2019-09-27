package wraperr

import (
	"fmt"
	"runtime"
	"strings"
)

func WithStack(err error, callArgs ...interface{}) error {
	return &withStack{
		err:   err,
		args:  callArgs,
		stack: callStack(1),
	}
}

func WithStackSkip(skip int, err error, callArgs ...interface{}) error {
	return &withStack{
		err:   err,
		args:  callArgs,
		stack: callStack(1 + skip),
	}
}

type withStack struct {
	err   error
	args  []interface{}
	stack []uintptr
}

func (w *withStack) Error() string {
	var b strings.Builder
	fmt.Fprintln(&b, w.err.Error())
	frames := runtime.CallersFrames(w.stack)
	frame, ok := frames.Next()
	if !ok {
		panic("no stack frame")
	}
	fmt.Fprintf(&b, "%s(%s)\n\t%s:%d", frame.Func.Name(), strings.Join(formatArgs(w.args), ", "), frame.File, frame.Line)
	for ok {
		fmt.Fprintf(&b, "%s\n\t%s:%d", frame.Func.Name(), frame.File, frame.Line)
		frame, ok = frames.Next()
	}
	return b.String()
}

func (w *withStack) Unwrap() error {
	return w.err
}

func callStack(skip int) []uintptr {
	c := make([]uintptr, 32)
	n := runtime.Callers(skip+1, c)
	return c[:n]
}
