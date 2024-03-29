package repo_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cypherium/cypherscan-server/internal/repo"
)

func TestBigIntWithJson(t *testing.T) {
	type A struct {
		BigInt repo.BigInt `json:"bigInt"`
	}

	a := A{repo.BigInt(*big.NewInt(98765432101))}

	b, err := json.Marshal(a)
	assert.NoError(t, err)
	assert.Equal(t, "{\"bigInt\":\"98765432101\"}", string(b))

	a1 := A{}
	err = json.Unmarshal(b, &a1)
	assert.NoError(t, err)
	assert.Equal(t, a, a1)
}

func TestHashWithJson(t *testing.T) {
	type A struct {
		Hash repo.Hash `json:"hash"`
	}

	a := A{repo.Hash{0x34, 0x33, 0x32, 0x31, 0x34, 0x33, 0x32, 0x31, 0x34, 0x33, 0x32, 0x31, 0x34, 0x33, 0x32, 0x31, 0x34, 0x33, 0x32, 0x31, 0x34, 0x33, 0x32, 0x31, 0x34, 0x33, 0x32, 0x31, 0x34, 0x33, 0x32, 0x31}}

	b, err := json.Marshal(a)
	assert.NoError(t, err)
	assert.Equal(t, "{\"hash\":\"0x3433323134333231343332313433323134333231343332313433323134333231\"}", string(b))

	a1 := A{}
	err = json.Unmarshal(b, &a1)
	assert.NoError(t, err)
	assert.Equal(t, a, a1)
}
func TestBytesWithJson(t *testing.T) {
	type A struct {
		Bytes repo.Bytes `json:"bytes"`
	}

	a := A{repo.Bytes{1, 2}}

	b, err := json.Marshal(a)
	assert.NoError(t, err)
	assert.Equal(t, "{\"bytes\":\"0x0102\"}", string(b))

	a1 := A{}
	err = json.Unmarshal(b, &a1)
	assert.NoError(t, err)
	assert.Equal(t, a, a1)
}
