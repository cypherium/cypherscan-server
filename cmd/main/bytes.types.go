package main

import (
	"encoding/hex"
	"fmt"
)

// Bytes is []byte
type Bytes []byte

// MarshalJSON is to support json
func (role Bytes) MarshalJSON() ([]byte, error) {
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
