package errors

// String implements the error interface for a string
// so it can be used as const error instead of var.
type String string

// Error implements the error interface
func (s String) Error() string {
	return string(s)
}
