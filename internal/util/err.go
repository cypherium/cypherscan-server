package util

import "fmt"

// MyError is the customised error type
type MyError struct {
	Message string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}
