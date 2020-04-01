package repo

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"

	"github.com/cypherium/cypherBFT/go-cypherium/common"
)

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

func (role Address) String() string {
	r := common.Address(role)
	return r.Hex()
}
