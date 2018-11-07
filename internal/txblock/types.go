package txblock

import (
  "database/sql/driver"
  "encoding/binary"
  "math/big"
)

// BigInt is big.Int
type BigInt big.Int

// Value is the Sacn interface
func (role BigInt) Value() (driver.Value, error) {
  return (*big.Int)(&role).Bytes(), nil
}

// Scan is the Scan interface
func (role *BigInt) Scan(value interface{}) error {
  (*big.Int)(role).SetBytes(value.([]byte))
  return nil
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
