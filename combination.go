package errors

import (
	"fmt"
	"io"
	"strings"
)

const multiErrorSeparator = "\n"

type multiError interface {
	Errors() []error
}

// combination combines multiple errors into one.
// The Error method returns the strings from the individual Error methods
// joined by the new line character '\n'.
// combination always has at least one error and returns the first error
// as result of the Cause method.
type combination struct {
	errs []error
	*stack
}

// Combine returns a combination error for 2 or more errors which are not nil,
// or the the single error if only one error was passed,
// or nil if zero arguments are passed or all passed errors are nil.
// The returned non nil error will be wrapped with the callstack,
// of the Combine call.
// The Combination type's Error method returns the strings from the
// individual Error methods joined by the new line character '\n'.
// Note that Cause(error) can only return a single error,
// so in case of a combination error, Cause returns the cause of the first error.
// When a passed error is a combination error implementing the following interface:
//     interface {
//         Errors() []error
//     }
// then the errors are flattened to form a new combination error together
// with the other passed errors.
func Combine(errs ...error) error {
	flattened := flatten(errs)

	switch len(flattened) {
	case 0:
		return nil
	case 1:
		return &withStack{
			flattened[0],
			callers(0),
		}
	}

	return &combination{
		flattened,
		callers(0),
	}
}

func flatten(errs []error) []error {
	var flattened []error
	for _, err := range errs {
		switch x := err.(type) {
		case nil:
			// ignore
		case multiError:
			flattened = append(flattened, x.Errors()...)
		default:
			flattened = append(flattened, x)
		}
	}
	return flattened
}

// Uncombine returns multible errors
// if err is a combination of multiple errors
// detected by implementing the following interface:
//     interface {
//         Errors() []error
//     }
// It returns the passed error in a single element slice
// if that error was not an error combination,
// or nil if the passed error was nil.
func Uncombine(err error) []error {
	if err == nil {
		return nil
	}
	if multi, ok := err.(multiError); ok {
		return multi.Errors()
	}
	return []error{err}
}

func (c *combination) Error() string {
	var b strings.Builder
	for i, err := range c.errs {
		if i > 0 {
			b.WriteString(multiErrorSeparator)
		}
		b.WriteString(err.Error())
	}
	return b.String()
}

func (c *combination) Cause() error {
	if len(c.errs) == 0 {
		return nil
	}
	return Cause(c.errs[0])
}

func (c *combination) Unwrap() error {
	return c.Cause()
}

func (c *combination) Errors() []error {
	return c.errs
}

func (c *combination) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			for _, e := range c.errs {
				fmt.Fprintf(s, "%+v\n", Cause(e).Error())
			}
			c.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, c.Error())
	case 'q':
		fmt.Fprintf(s, "%q", c.Error())
	}
}
