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
    t       ConvertedType
    out     interface{}
    isError bool
  }{
    {"345", BigIntType, *big.NewInt(0x345), false},
    {"0xabc", BigIntType, *big.NewInt(0xabc), false},
    {"0Xbc", BigIntType, *big.NewInt(0xbc), false},
    {"0xzzd", BigIntType, *big.NewInt(0), true},
    {"0x1122", HashType, &Hash{0x11, 0x22}, false},
    {"0x3344", AddressType, &Address{0x33, 0x44}, false},
    {"0x4455", BloomType, &Bloom{0x44, 0x55}, false},
    {"0x5566", BlockNonceType, &BlockNonce{0x55, 0x66}, false},
  }

  for _, table := range tables {
    if !table.isError {
      out := Parse(table.in, table.t)
      assert.Equal(t, table.out, out)
    }
  }
}
