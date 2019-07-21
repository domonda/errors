package errors

// Const implements the error interface for a string.
// Use this type to declare package level errors as const instead of var.
// Example:
//     const ErrThisIsConst = errors.Const("this is const")
type Const string

// Error implements the error interface
func (s Const) Error() string {
	return string(s)
}
