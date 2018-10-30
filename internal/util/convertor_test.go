package util

import (
  "math/big"
  "testing"
)

func TestStrip0x(t *testing.T) {
  tables := []struct {
    in  string
    out string
  }{
    {"345", "345"},
    {"0xabc", "abc"},
    {"0Xbc", "bc"},
  }

  for _, table := range tables {
    out := Stripe0x(table.in)
    if out != table.out {
      t.Errorf("Strip0x (%s) was incorrect, got: %s, expect: %s.", table.in, out, table.out)
    }
  }
}

func TestHxStrToBigInt(t *testing.T) {
  tables := []struct {
    in      string
    out     *big.Int
    isError bool
  }{
    {"345", big.NewInt(0x345), false},
    {"0xabc", big.NewInt(0xabc), false},
    {"0Xbc", big.NewInt(0xbc), false},
    {"0xzzd", big.NewInt(0), true},
  }

  for _, table := range tables {
    out, err := HxStrToBigInt(table.in)
    if out.Cmp(table.out) != 0 || (table.isError && err == nil) {
      t.Errorf("HxStrToBigInt (%s) was incorrect, got: %s, expect: %s.", table.in, out, table.out)
    }
  }
}
