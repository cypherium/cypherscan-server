package util

import (
  "encoding/hex"
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

// HxStrToHash convert hx string to Hash
func HxStrToHash(s string) (*Hash, error) {
  var ret Hash
  r := strings.NewReader(Stripe0x(s))
  reader := hex.NewDecoder(r)
  _, err := reader.Read(ret[:])
  return &ret, err
}

// HxStrToAddress convert hx string to Address
func HxStrToAddress(s string) (*Address, error) {
  var ret Address
  r := strings.NewReader(Stripe0x(s))
  reader := hex.NewDecoder(r)
  _, err := reader.Read(ret[:])
  return &ret, err
}

// HxStrToBloom convert hx string to Address
func HxStrToBloom(s string) (*Bloom, error) {
  var ret Bloom
  r := strings.NewReader(Stripe0x(s))
  reader := hex.NewDecoder(r)
  _, err := reader.Read(ret[:])
  return &ret, err
}

// ConvertedType is a enum
type ConvertedType int

const (
  // BigIntType is *big.Int
  BigIntType ConvertedType = 1 + iota
  // HashType is Hash
  HashType
  // AddressType is Address
  AddressType
  // BloomType is Bloom
  BloomType
)

// Parse is a generic convert function will take care of error
func Parse(in string, t ConvertedType) interface{} {
  switch t {
  case BigIntType:
    out, err := HxStrToBigInt(in)
    if err != nil {
      log.Println("convert error from Parse")
    }
    return out
  case HashType:
    out, err := HxStrToHash(in)
    if err != nil {
      log.Println("convert error from Parse")
    }
    return out
  case AddressType:
    out, err := HxStrToAddress(in)
    if err != nil {
      log.Println("convert error from Parse")
    }
    return out
  case BloomType:
    out, err := HxStrToBloom(in)
    if err != nil {
      log.Println("convert error from Parse")
    }
    return out
  default:
    return in
  }
}
