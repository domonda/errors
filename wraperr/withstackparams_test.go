package wraperr

import (
	"fmt"
	"runtime"
	"strings"
)

func funcA(i int, s string) (err error) {
	defer WithCallParams(&err, i, s)

	return funcB(s)
}

func funcB(s ...string) (err error) {
	defer WithCallParams(&err, s)

	return funcC()
}

func funcC() (err error) {
	defer WithCallParams(&err)

	return New("error in funcC")
}

func basePath() string {
	_, file, _, _ := runtime.Caller(1)
	return file[:strings.Index(file, "github.com")]
}

func ExampleWithCallParams() {
	err := funcA(666, "Hello World!")
	errStr := err.Error()
	errStr = strings.ReplaceAll(errStr, basePath(), "")
	fmt.Println(errStr)

	// Output:
	// error in funcC
	// github.com/domonda/errors/wraperr.funcC()
	//     github.com/domonda/errors/wraperr/withstackparams_test.go:24
	// github.com/domonda/errors/wraperr.funcB([]string{"Hello World!"})
	//     github.com/domonda/errors/wraperr/withstackparams_test.go:18
	// github.com/domonda/errors/wraperr.funcA(666, "Hello World!")
	//     github.com/domonda/errors/wraperr/withstackparams_test.go:12
}
