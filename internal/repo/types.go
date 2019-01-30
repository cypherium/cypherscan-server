package repo

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/cypherium/CypherTestNet/go-cypherium/common"
	"gitlab.com/ron-liu/cypherscan-server/internal/util"
)

// BigInt is big.Int
type BigInt big.Int

// Value is the Sacn interface
func (i BigInt) Value() (driver.Value, error) {
	// return (*big.Int)(&i).Bytes(), nil
	return (*big.Int)(&i).String(), nil
}

// Scan is the Scan interface
func (i *BigInt) Scan(value interface{}) error {
	// (*big.Int)(i).SetBytes(value.([]byte))
	(*big.Int)(i).SetString(value.(string), 10)
	return nil
}

// MarshalJSON is to support json
func (i BigInt) MarshalJSON() ([]byte, error) {
	i2 := big.Int(i)
	return []byte(fmt.Sprintf(`"%s"`, i2.String())), nil
}

// UnmarshalJSON is to support json
func (i *BigInt) UnmarshalJSON(b []byte) error {
	z := (*big.Int)(i)
	s := string(b[1 : len(b)-1])
	_, ok := z.SetString(s, 10)
	if !ok {
		return &util.MyError{Message: fmt.Sprintf("Error to Unmarshal to big.Int: %s", s)}
	}
	return nil
}

// UInt64 is uint64
type UInt64 uint64

// Value is the Sacn interface
func (role UInt64) Value() (driver.Value, error) {
	x := int64(role)
	return x, nil
}

// Scan is the Scan interface
func (role *UInt64) Scan(value interface{}) error {
	x := value.(int64)
	*role = UInt64(x)

	return nil
}

// Hash is common.Hash
type Hash common.Hash

// Value is the Sacn interface
func (role Hash) Value() (driver.Value, error) {
	b := fmt.Sprintf("%x", role)
	return b, nil
}

// Scan is the Scan interface
func (role *Hash) Scan(value interface{}) error {
	data, err := hex.DecodeString(string(value.([]byte)))
	if err != nil {
		return err
	}
	r := (*common.Hash)(role)
	copy(r[:], data)
	return nil
}

// Hex is to change the Hash to hex string
func (role Hash) Hex() string {
	r := common.Hash(role)
	return r.Hex()
}

// MarshalJSON is to support json
func (role Hash) MarshalJSON() ([]byte, error) {
	dst := make([]byte, hex.EncodedLen(len(role)))
	hex.Encode(dst, role[:])
	return []byte(fmt.Sprintf(`"0x%s"`, dst)), nil
}

// UnmarshalJSON is to support json
func (role *Hash) UnmarshalJSON(b []byte) error {
	src := string(b[3 : len(b)-1])
	bytes, err := hex.DecodeString(src)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return err
	}
	for i, b := range bytes {
		role[i] = b
	}
	return nil
}

// Address is common.Hash
type Address common.Address

// MarshalJSON is to support json
func (role Address) MarshalJSON() ([]byte, error) {
	dst := make([]byte, hex.EncodedLen(len(role)))
	hex.Encode(dst, role[:])
	return []byte(fmt.Sprintf(`"0x%s"`, dst)), nil
}

// UnmarshalJSON is to support json
func (role *Address) UnmarshalJSON(b []byte) error {
	src := string(b[3 : len(b)-1])
	bytes, err := hex.DecodeString(src)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return err
	}
	for i, b := range bytes {
		role[i] = b
	}
	return nil
}

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
