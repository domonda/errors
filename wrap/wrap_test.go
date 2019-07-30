package wrap

import (
	"errors"
	"testing"

	"github.com/domonda/go-types/uu"
)

func Test_formatResultError(t *testing.T) {
	var (
		str    string
		strPtr *string
		uid    uu.ID
		uidPtr *uu.ID
		i      int
		iPtr   *int
		interf interface{}
	)

	result := FormatCallSignature("test", str, strPtr, uid, uidPtr, i, iPtr, interf)
	expected := `test("", <nil>, "00000000-0000-0000-0000-000000000000", <nil>, 0, <nil>, <nil>)`
	if result != expected {
		t.Errorf("result `%s` != expected `%s`", result, expected)
	}
}

func wrappedErrorFunc() (err error) {
	defer ResultError(&err, "wrappedErrorFunc")

	return errors.New("TEST")
}

func Test_StackFrame(t *testing.T) {
	err := wrappedErrorFunc()
	t.Logf("%+v", err)
	// t.Fail()
}
