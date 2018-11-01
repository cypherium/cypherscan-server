package util

import (
  "fmt"
  "log"
  "math/big"
  "strings"
)

// MyError is the customised error type
type MyError struct {
  message string
}

func (e *MyError) Error() string {
  return fmt.Sprintf("%s", e.message)
}

// Stripe0x stripe the prefixed 0x if existed, otherwise return unchanged
func Stripe0x(s string) string {
  s = strings.TrimPrefix(s, "0x")
  return strings.TrimPrefix(s, "0X")
}

// HxStrToBigInt convert hx string like "0xff" to big.Int
func HxStrToBigInt(s string) (n *big.Int, err error) {
  striped := Stripe0x(s)

  n = new(big.Int)
  _, ok := n.SetString(striped, 16)
  if !ok {
    err = &MyError{fmt.Sprintf("failed when convert (%s) to big.Int", s)}
  }
  return
}

// Parse is a generic convert function will take care of error
func Parse(in string, sample interface{}) interface{} {
  switch sample.(type) {
  case *big.Int:
    out, err := HxStrToBigInt(in)
    if err != nil {
      log.Println("convert error from Parse")
    }
    return out
  default:
    return in
  }
}
