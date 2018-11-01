package util

import (
  "github.com/stretchr/testify/assert"
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
    assert.Equal(t, out, table.out, "Strip0x (%s) was incorrect, got: %s, expect: %s.", table.in, out, table.out)
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
    assert.True(t, out.Cmp(table.out) == 0 || (table.isError && err != nil), "HxStrToBigInt (%s) was incorrect, got: %s, expect: %s.", table.in, out, table.out)
  }
}

func TestParse(t *testing.T) {
  tables := []struct {
    in      string
    sample  interface{}
    out     *big.Int
    isError bool
  }{
    {"345", new(big.Int), big.NewInt(0x345), false},
    {"0xabc", new(big.Int), big.NewInt(0xabc), false},
    {"0Xbc", new(big.Int), big.NewInt(0xbc), false},
    {"0xzzd", new(big.Int), big.NewInt(0), true},
  }

  for _, table := range tables {
    if !table.isError {
      out := Parse(table.in, table.sample).(*big.Int)
      assert.True(t, out.Cmp(table.out) == 0, "Parse (%s) was incorrect, got: %s, expect: %s.", table.in, out, table.out)
    }

  }
}
