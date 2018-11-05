package util

import (
  "github.com/stretchr/testify/assert"
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

func TestParse(t *testing.T) {
  tables := []struct {
    in      string
    t       ConvertedType
    out     interface{}
    isError bool
  }{
    {"0x1122", BytesType, []byte{0x11, 0x22}, false},
    {"0x3344", BytesType, []byte{0x33, 0x44}, false},
    {"0x4455", BytesType, []byte{0x44, 0x55}, false},
    {"0x5566", BytesType, []byte{0x55, 0x66}, false},
    {"0x556", BytesType, []byte{0x05, 0x56}, false},
    {"0x556", UInt64Type, uint64(0x556), false},
    {"", UInt64Type, uint64(0), false},
    {"0x516", UInt32Type, uint32(0x516), false},
    {"", UInt32Type, uint32(0), false},
  }

  for _, table := range tables {
    if !table.isError {
      out := Parse(table.in, table.t)
      assert.Equal(t, table.out, out, table.in)
    }
  }
}
