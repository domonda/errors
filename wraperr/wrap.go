package wraperr

import (
	"errors"
	"fmt"
	"runtime"
)

func callStack(skip int) []uintptr {
	c := make([]uintptr, 32)
	n := runtime.Callers(skip+2, c)
	return c[:n]
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

func Result(resultVar *error, params ...interface{}) {
	if *resultVar == nil {
		return
	}
	*resultVar = &withStackParams{
		withStack: withStack{
			err:   *resultVar,
			stack: callStack(1),
		},
		params: params,
	}
}

func Result0(resultVar *error) {
	if *resultVar == nil {
		return
	}
	*resultVar = &withStackParams{
		withStack: withStack{
			err:   *resultVar,
			stack: callStack(1),
		},
		params: nil,
	}
}

func Result1(resultVar *error, p0 interface{}) {
	if *resultVar == nil {
		return
	}
	*resultVar = &withStackParams{
		withStack: withStack{
			err:   *resultVar,
			stack: callStack(1),
		},
		params: []interface{}{p0},
	}
}

func Result2(resultVar *error, p0, p1 interface{}) {
	if *resultVar == nil {
		return
	}
	*resultVar = &withStackParams{
		withStack: withStack{
			err:   *resultVar,
			stack: callStack(1),
		},
		params: []interface{}{p0, p1},
	}
}

func Result3(resultVar *error, p0, p1, p2 interface{}) {
	if *resultVar == nil {
		return
	}
	*resultVar = &withStackParams{
		withStack: withStack{
			err:   *resultVar,
			stack: callStack(1),
		},
		params: []interface{}{p0, p1, p2},
	}
}
