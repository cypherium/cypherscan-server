package repo

import (
	"fmt"
)

func bytesToPostgresSearchableString(bytes []byte) string {
	return fmt.Sprintf("\\x%x", bytes)
}
