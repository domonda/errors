package wraperr

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
