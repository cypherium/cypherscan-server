package util

import (
  "encoding/hex"
  "fmt"
  log "github.com/sirupsen/logrus"
  // "math/big"
  "math/big"
  "strconv"
  "strings"
  "time"
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

// HxStrToBytes convert hx string to []byte
func HxStrToBytes(s string) ([]byte, error) {
  if s == "" {
    return make([]byte, 0), nil
  }
  striped := Stripe0x(s)

  return hex.DecodeString(
    func(x string) string {
      if len(x)%2 > 0 {
        return "0" + x
      }
      return x
    }(striped))
}

// HxStrToUInt32 convert hx string to uint32
func HxStrToUInt32(s string) (uint32, error) {
  if s == "" {
    return 0, nil
  }
  n, err := strconv.ParseUint(Stripe0x(s), 16, 32)
  return uint32(n), err
}

// HxStrToUInt64 convert hx string to uint64
func HxStrToUInt64(s string) (uint64, error) {
  if s == "" {
    return 0, nil
  }
  return strconv.ParseUint(Stripe0x(s), 16, 64)
}

// HxStrToTime convert hx string to time
func HxStrToTime(s string) (time.Time, error) {
  i, err := strconv.ParseInt(Stripe0x(s), 16, 64)
  if err != nil {
    return time.Time{}, err
  }
  tm := time.Unix(i, 0)
  return tm, err
}

// ConvertedType is a enum
type ConvertedType int

const (
  // UInt64Type is unit64
  UInt64Type = 1 + iota
  // TimeType is Time
  TimeType
  // BytesType is []byte
  BytesType
  // UInt32Type is unit
  UInt32Type
)

// Parse is a generic convert function will take care of error
func Parse(in string, t ConvertedType) interface{} {
  switch t {
  case BytesType:
    out, err := HxStrToBytes(in)
    if err != nil {
      log.Printf("convert error from Parse: %v %v %v", in, t, err)
    }
    return out
  case UInt64Type:
    out, err := HxStrToUInt64(in)
    if err != nil {
      log.Println("convert error from Parse:", t, err)
    }
    return out
  case UInt32Type:
    out, err := HxStrToUInt32(in)
    if err != nil {
      log.Println("convert error from Parse:", t, err)
    }
    return out
  case TimeType:
    out, err := HxStrToTime(in)
    if err != nil {
      log.Println("convert error from Parse:", t, err)
    }
    return out
  default:
    return in
  }
}

// StringToBigInt is to change string to big.Int
func StringToBigInt(s string, base int) (*big.Int, error) {
  n := new(big.Int)
  n, ok := n.SetString(s, base)
  if !ok {
    log.Error("StringToBigInt error")
    return n, &MyError{"StringToBigInt error"}
  }
  return n, nil
}
