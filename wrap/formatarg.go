package wrap

import (
	"encoding/json"
	"fmt"
	"reflect"

	reflection "github.com/ungerik/go-reflection"
)

func formatArg(arg interface{}) string {
	if arg == nil {
		return "<nil>"
	}
	v := reflect.ValueOf(arg)
	if reflection.IsNil(v) {
		return "<nil>"
	}

	switch a := arg.(type) {
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

	switch t := reflection.DerefType(v.Type()); t.Kind() {
	case reflect.Func:
		return "<func>"
	case reflect.Chan:
		return "<chan>"
	case reflect.Struct:
		bytes, err := json.Marshal(arg)
		if err != nil {
			return t.Name() + "marshaling error"
		}
		return t.Name() + string(bytes)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// "%#v" would return hex literal
		return fmt.Sprintf("%v", arg)
	}

	return fmt.Sprintf("%#v", arg)
}
