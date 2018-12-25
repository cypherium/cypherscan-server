package repo

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/cypherium/CypherTestNet/go-cypherium/common"
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

// Hash is common.Hash
type Hash common.Hash

// Value is the Sacn interface
func (role Hash) Value() (driver.Value, error) {
	b := role[:]
	return b, nil
}

// Scan is the Scan interface
func (role *Hash) Scan(value interface{}) error {
	r := (*common.Hash)(role)
	copy(r[:], value.([]byte))
	return nil
}

// Hex is to change the Hash to hex string
func (role Hash) Hex() string {
	r := common.Hash(role)
	return r.Hex()
}

// Address is common.Hash
type Address common.Address

// Value is the Sacn interface
func (role Address) Value() (driver.Value, error) {
	b := role[:]
	return b, nil
}

// Scan is the Scan interface
func (role *Address) Scan(value interface{}) error {
	r := (*common.Address)(role)
	copy(r[:], value.([]byte))
	return nil
}

// Hex is to change the Hash to hex string
func (role Address) Hex() string {
	r := common.Address(role)
	return r.Hex()
}
