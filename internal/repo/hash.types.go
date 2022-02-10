package repo

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"

	"github.com/cypherium/cypherBFT/common"
)

// Lengths of hashes and addresses in bytes.
const (
	HashLength    = 32
	AddressLength = 20
	PubkeyLength  = 32
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

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// Hex is to change the Hash to hex []byte
func (role Hash) Bytes() []byte {
	r := common.Hash(role)
	return r.Bytes()
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
