package txblock

import (
  "database/sql/driver"
  "encoding/binary"
  "fmt"
  "math/big"
)

// BigInt is big.Int
type BigInt big.Int

// Value is the Sacn interface
func (i BigInt) Value() (driver.Value, error) {
  return (*big.Int)(&i).Bytes(), nil
}

// Scan is the Scan interface
func (i *BigInt) Scan(value interface{}) error {
  (*big.Int)(i).SetBytes(value.([]byte))
  return nil
}

// MarshalJSON is to support json
func (i BigInt) MarshalJSON() ([]byte, error) {
  i2 := big.Int(i)
  return []byte(fmt.Sprintf(`"%s"`, i2.String())), nil
}

// UInt64 is uint64
type UInt64 uint64

// Value is the Sacn interface
func (role UInt64) Value() (driver.Value, error) {
  b := make([]byte, 8)
  binary.LittleEndian.PutUint64(b, uint64(role))
  return b, nil
}

// Scan is the Scan interface
func (role *UInt64) Scan(value interface{}) error {
  *role = UInt64(binary.LittleEndian.Uint64(value.([]byte)))
  return nil
}
