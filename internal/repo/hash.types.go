package repo

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"

	"github.com/cypherium/cypherBFT/go-cypherium/common"
)

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
func (role Hash) String() string {
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
