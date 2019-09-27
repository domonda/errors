package wraperr_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/domonda/errors/wraperr"
)

func funcA(i int, s string) (err error) {
	defer wraperr.ResultVar(&err, i, s)

	return funcB(s)
}

func funcB(s ...string) (err error) {
	defer wraperr.ResultVar(&err, s)

	return funcC()
}

func funcC() (err error) {
	defer wraperr.ResultVar(&err)

	return wraperr.New("error in funcC")
}

func Test_ResultVar(t *testing.T) {
	err := funcA(666, "Hello World!")
	assert.Error(t, err, "test error")

	fmt.Println(err)
	t.Fail() // to see Println above
}
