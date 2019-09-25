package wrap

import (
	stderrors "errors"
	"fmt"
	"strings"
	"time"

	"github.com/domonda/errors"
)

type Logger interface {
	Printf(format string, args ...interface{})
}

func Error(err error, funcName string, funcArgs ...interface{}) error {
	if err == nil {
		return nil
	}

	return errors.WrapSkip(1, err, callSignature(funcName, funcArgs))
}

func ResultError(errPtr *error, funcName string, funcArgs ...interface{}) {
	if *errPtr == nil {
		return
	}

	*errPtr = errors.WrapSkip(1, *errPtr, callSignature(funcName, funcArgs))
}

func RecoverPanicAsResultError(errPtr *error, funcName string, funcArgs ...interface{}) {
	err := AsError(recover())
	if err == nil {
		return
	}

	err = errors.WrapSkip(1, err, callSignature(funcName, funcArgs))

	*errPtr = errors.Combine(err, *errPtr)
}

func LogPanic(log Logger, funcName string, funcArgs ...interface{}) {
	p := recover()
	if p == nil {
		return
	}

	err := errors.WrapSkip(1, AsError(p), callSignature(funcName, funcArgs))

	log.Printf("LogPanic: %w", err)

	panic(p)
}

func RecoverAndLogPanic(log Logger, funcName string, funcArgs ...interface{}) {
	err := AsError(recover())
	if err == nil {
		return
	}

	err = errors.WrapSkip(1, err, callSignature(funcName, funcArgs))

	log.Printf("RecoverAndLogPanic: %w", err)
}

func AsError(val interface{}) error {
	switch x := val.(type) {
	case nil:
		return nil
	case error:
		return x
	case string:
		return stderrors.New(x)
	case fmt.Stringer:
		return stderrors.New(x.String())
	}
	return stderrors.New(fmt.Sprintf("%+v", val))
}

func callSignature(funcName string, funcArgs []interface{}) string {
	var b strings.Builder
	b.WriteString("CALL: ")
	return formatCallSignature(&b, funcName, funcArgs)
}

func formatCallSignature(b *strings.Builder, funcName string, funcArgs []interface{}) string {
	b.WriteString(funcName)
	b.WriteByte('(')
	for i, arg := range funcArgs {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(formatArg(arg))
	}
	b.WriteByte(')')

	return b.String()
}

func FormatCallSignature(funcName string, funcArgs ...interface{}) string {
	var b strings.Builder
	return formatCallSignature(&b, funcName, funcArgs)
}

func DerefString(s *string) interface{} {
	if s == nil {
		return nil
	}
	return *s
}

func DerefTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return *t
}
