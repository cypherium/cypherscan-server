package repo_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
)

func TestCustomizedTypeWithJson(t *testing.T) {
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
