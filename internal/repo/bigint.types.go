package repo

import (
	"database/sql/driver"
	"fmt"
	"math/big"

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

// MarshalJSON is to support json
func (role UInt64) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"0x%x"`, uint64(role))), nil
}
