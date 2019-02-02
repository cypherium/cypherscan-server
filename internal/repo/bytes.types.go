package repo

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
)

// Bytes is []byte
type Bytes []byte

// Value is the Sacn interface
func (role Bytes) Value() (driver.Value, error) {
	b := ([]byte)(role)
	return b, nil
}

// Scan is the Scan interface
func (role *Bytes) Scan(value interface{}) error {
	r := value.([]byte)
	*role = r
	// copy(r[:], value.([]byte))
	return nil
}

// MarshalJSON is to support json
func (role Bytes) MarshalJSON() ([]byte, error) {
	// dst := make([]byte, hex.EncodedLen(len(role)))
	// hex.Encode(dst, role[:])
	// return []byte(fmt.Sprintf(`"0x%s"`, dst)), nil
	return []byte(fmt.Sprintf(`"0x%x"`, []byte(role))), nil
}

// UnmarshalJSON is to support json
func (role *Bytes) UnmarshalJSON(b []byte) error {
	src := string(b[3 : len(b)-1])
	bytes, err := hex.DecodeString(src)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return err
	}
	*role = bytes
	return nil
}

func (role Bytes) String() string {
	return fmt.Sprintf("0x%x", ([]byte)(role))
}
