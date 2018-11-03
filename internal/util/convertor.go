package util

import (
  "encoding/hex"
  "fmt"
  "log"
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

// HxStrToBigInt convert hx string like "0xff" to big.Int
func HxStrToBigInt(s string) (n big.Int, err error) {
  striped := Stripe0x(s)
  n = *new(big.Int)
  if s == "" {
    return n, nil
  }

  _, ok := n.SetString(striped, 16)
  if !ok {
    err = &MyError{fmt.Sprintf("failed when convert (%s) to big.Int", s)}
  }
  return
}

func hxStrToBytes(s string, bytes []byte) ([]byte, error) {
  if s == "" {
    return bytes, nil
  }
  r := strings.NewReader(Stripe0x(s))
  reader := hex.NewDecoder(r)
  _, err := reader.Read(bytes)
  return bytes, err
}

// HxStrToHash convert hx string to Hash
func HxStrToHash(s string) (*Hash, error) {
  var ret Hash
  _, err := hxStrToBytes(s, ret[:])
  return &ret, err
}

// HxStrToAddress convert hx string to Address
func HxStrToAddress(s string) (*Address, error) {
  var ret Address
  _, err := hxStrToBytes(s, ret[:])
  return &ret, err
}

// HxStrToBloom convert hx string to Address
func HxStrToBloom(s string) (*Bloom, error) {
  var ret Bloom
  _, err := hxStrToBytes(s, ret[:])
  return &ret, err
}

// HxStrToBlockNonce convert hx string to BlockNonce
func HxStrToBlockNonce(s string) (*BlockNonce, error) {
  var ret BlockNonce
  r := strings.NewReader(Stripe0x(s))
  reader := hex.NewDecoder(r)
  _, err := reader.Read(ret[:])
  return &ret, err
}

// HxStrToUInt64 convert hx string to uint64
func HxStrToUInt64(s string) (uint64, error) {
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
  // BigIntType is *big.Int
  BigIntType ConvertedType = 1 + iota
  // HashType is Hash
  HashType
  // AddressType is Address
  AddressType
  // BloomType is Bloom
  BloomType
  // BlockNonceType is BlockNonce
  BlockNonceType
  // UInt64Type is unit64
  UInt64Type
  // TimeType is Time
  TimeType
)

// Parse is a generic convert function will take care of error
func Parse(in string, t ConvertedType) interface{} {
  switch t {
  case BigIntType:
    out, err := HxStrToBigInt(in)
    if err != nil {
      log.Println("convert error from Parse:", err)
    }
    return out
  case HashType:
    out, err := HxStrToHash(in)
    if err != nil {
      log.Println("convert error from Parse:", err)
    }
    return out
  case AddressType:
    out, err := HxStrToAddress(in)
    if err != nil {
      log.Println("convert error from Parse:", err)
    }
    return out
  case BloomType:
    out, err := HxStrToBloom(in)
    if err != nil {
      log.Println("convert error from Parse:", err)
    }
    return out
  case BlockNonceType:
    out, err := HxStrToBlockNonce(in)
    if err != nil {
      log.Println("convert error from Parse:", t)
    }
    return out
  case UInt64Type:
    out, err := HxStrToUInt64(in)
    if err != nil {
      log.Println("convert error from Parse:", t)
    }
    return out
  case TimeType:
    out, err := HxStrToTime(in)
    if err != nil {
      log.Println("convert error from Parse:", t)
    }
    return out
  default:
    return in
  }
}
