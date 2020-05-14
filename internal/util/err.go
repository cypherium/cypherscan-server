package util

import (
	"fmt"
	"runtime/debug"
)

// MyError is the customised error type
type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

// LowLevelError is the low level error
type LowLevelError struct {
	error
}

func (e *MyError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

// NewError will create an error
func NewError(err error, message string, messageArgs ...interface{}) *MyError {
	return &MyError{
		Inner:      err,
		Message:    fmt.Sprintf(message, messageArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}),
	}
}

// PickFirstFromErrs is to pick the first error from errors
func PickFirstFromErrs(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
