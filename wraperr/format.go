package wraperr

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	reflection "github.com/ungerik/go-reflection"
)

func formatError(err error) string {
	var (
		calls             []string
		firstWithoutStack error
	)

	for err != nil {
		switch e := err.(type) {
		case callStackParamsProvider:
			calls = append(calls, formatCallParamsSource(e.CallStackParams()))

		case callStackProvider:
			calls = append(calls, formatCallSource(e.CallStack()))

		default:
			if firstWithoutStack == nil {
				firstWithoutStack = err
			}
		}

		err = errors.Unwrap(err)
	}

	var b strings.Builder
	fmt.Fprintln(&b, firstWithoutStack.Error())
	for i := len(calls) - 1; i >= 0; i-- {
		fmt.Fprintln(&b, calls[i])
	}
	return b.String()
}

// func formatCall(stack []uintptr) string {
// 	frame, ok := runtime.CallersFrames(stack).Next()
// 	if !ok {
// 		return "insufficient call stack"
// 	}
// 	return fmt.Sprintf("%s", frame.Function)
// }

func formatCallSource(stack []uintptr) string {
	frame, ok := runtime.CallersFrames(stack).Next()
	if !ok {
		return "insufficient call stack"
	}
	return fmt.Sprintf("%s\n    %s:%d", frame.Function, frame.File, frame.Line)
}

// func formatCallParams(stack []uintptr, params []interface{}) string {
// 	frame, ok := runtime.CallersFrames(stack).Next()
// 	if !ok {
// 		return "insufficient call stack"
// 	}
// 	return fmt.Sprintf("%s(%s)", frame.Function, formatParams(params))
// }

func formatCallParamsSource(stack []uintptr, params []interface{}) string {
	frame, ok := runtime.CallersFrames(stack).Next()
	if !ok {
		return "insufficient call stack"
	}
	return fmt.Sprintf("%s(%s)\n    %s:%d", frame.Function, formatParams(params), frame.File, frame.Line)
}

func formatParams(params []interface{}) string {
	var b strings.Builder
	for i, param := range params {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(formatParam(param))
	}
	return b.String()
}

func formatParam(param interface{}) string {
	if param == nil {
		return "<nil>"
	}
	v := reflect.ValueOf(param)
	if reflection.IsNil(v) {
		return "<nil>"
	}

	switch a := param.(type) {
	case error:
		return fmt.Sprintf("error(%q)", a.Error())
	case fmt.Stringer:
		return fmt.Sprintf("%q", a.String())
	case []byte:
		if len(a) > 300 {
			return fmt.Sprintf("[%d]byte(%q...)", len(a), a[:10])
		}
		return fmt.Sprintf("[]byte(%q)", a)
	}

	if v.Kind() == reflect.Ptr {
		switch v.Elem().Kind() {
		case reflect.Struct:
			// handle futher down

		case reflect.Func:
			return "<func>"

		case reflect.Chan:
			return "<chan>"

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// "%#v" would return hex literal
			return fmt.Sprintf("%v", v.Elem().Interface())

		default:
			return fmt.Sprintf("%#v", v.Elem().Interface())
		}
	}

	switch t := reflection.DerefType(v.Type()); t.Kind() {
	case reflect.Func:
		return "<func>"

	case reflect.Chan:
		return "<chan>"

	case reflect.Struct:
		bytes, err := json.Marshal(param)
		if err != nil {
			return t.Name() + "marshaling error"
		}
		return t.Name() + string(bytes)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// "%#v" would return hex literal
		return fmt.Sprintf("%v", param)
	}

	return fmt.Sprintf("%#v", param)
}
