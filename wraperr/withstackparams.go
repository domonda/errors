package wraperr

/*
Call argument parameters are available on the stack,
but in a platform dependent packed format and not directly accessible
via package runtime.
Could only be parsed from runtime.Stack() result text.

See https://www.ardanlabs.com/blog/2015/01/stack-traces-in-go.html
*/

func WithCallParams(resultVar *error, params ...interface{}) {
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

func With0CallParams(resultVar *error) {
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

func With1CallParams(resultVar *error, p0 interface{}) {
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

func With2CallParams(resultVar *error, p0, p1 interface{}) {
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

func With3CallParams(resultVar *error, p0, p1, p2 interface{}) {
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

func With4CallParams(resultVar *error, p0, p1, p2, p3 interface{}) {
	if *resultVar == nil {
		return
	}
	*resultVar = &withStackParams{
		withStack: withStack{
			err:   *resultVar,
			stack: callStack(1),
		},
		params: []interface{}{p0, p1, p2, p3},
	}
}

func With5CallParams(resultVar *error, p0, p1, p2, p3, p4 interface{}) {
	if *resultVar == nil {
		return
	}
	*resultVar = &withStackParams{
		withStack: withStack{
			err:   *resultVar,
			stack: callStack(1),
		},
		params: []interface{}{p0, p1, p2, p3, p4},
	}
}

func With6CallParams(resultVar *error, p0, p1, p2, p3, p4, p5 interface{}) {
	if *resultVar == nil {
		return
	}
	*resultVar = &withStackParams{
		withStack: withStack{
			err:   *resultVar,
			stack: callStack(1),
		},
		params: []interface{}{p0, p1, p2, p3, p4, p5},
	}
}

func With7CallParams(resultVar *error, p0, p1, p2, p3, p4, p5, p6 interface{}) {
	if *resultVar == nil {
		return
	}
	*resultVar = &withStackParams{
		withStack: withStack{
			err:   *resultVar,
			stack: callStack(1),
		},
		params: []interface{}{p0, p1, p2, p3, p4, p5, p6},
	}
}

func With8CallParams(resultVar *error, p0, p1, p2, p3, p4, p5, p6, p7 interface{}) {
	if *resultVar == nil {
		return
	}
	*resultVar = &withStackParams{
		withStack: withStack{
			err:   *resultVar,
			stack: callStack(1),
		},
		params: []interface{}{p0, p1, p2, p3, p4, p5, p6, p7},
	}
}

func With9CallParams(resultVar *error, p0, p1, p2, p3, p4, p5, p6, p7, p8 interface{}) {
	if *resultVar == nil {
		return
	}
	*resultVar = &withStackParams{
		withStack: withStack{
			err:   *resultVar,
			stack: callStack(1),
		},
		params: []interface{}{p0, p1, p2, p3, p4, p5, p6, p7, p8},
	}
}

func With10CallParams(resultVar *error, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9 interface{}) {
	if *resultVar == nil {
		return
	}
	*resultVar = &withStackParams{
		withStack: withStack{
			err:   *resultVar,
			stack: callStack(1),
		},
		params: []interface{}{p0, p1, p2, p3, p4, p5, p6, p7, p8, p9},
	}
}

type callStackParamsProvider interface {
	CallStackParams() ([]uintptr, []interface{})
}

type withStackParams struct {
	withStack

	params []interface{}
}

func (w *withStackParams) Error() string {
	return formatError(w)
}

func (w *withStackParams) CallStackParams() ([]uintptr, []interface{}) {
	return w.stack, w.params
}
