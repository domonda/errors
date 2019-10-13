package main

import (
	"errors"
	"fmt"

	"github.com/domonda/errors/wraperr"
)

func funcA(i int, s string) (err error) {
	defer wraperr.WithFuncParams(&err, i, s)

	return funcB(s)
}

func funcB(s ...string) (err error) {
	defer wraperr.WithFuncParams(&err, s)

	return funcC()
}

func funcC() (err error) {
	defer wraperr.WithFuncParams(&err)

	return errors.New("error in funcC")
}

func main() {
	err := funcA(666, "Hello World!")
	str := err.Error()
	fmt.Println(str)
}
