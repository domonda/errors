package wraperr_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/domonda/errors/wraperr"
)

func funcA(i int, s string) (err error) {
	defer wraperr.Result(&err, i, s)

	return funcB(s)
}

func funcB(s ...string) (err error) {
	defer wraperr.Result(&err, s)

	return funcC()
}

func funcC() (err error) {
	defer wraperr.Result(&err)

	return wraperr.New("error in funcC")
}

func Test_Resul(t *testing.T) {
	err := funcA(666, "Hello World!")
	assert.Error(t, err, "test error")

	str := err.Error()
	fmt.Println(str)
	t.Fail() // to see Println above
}
