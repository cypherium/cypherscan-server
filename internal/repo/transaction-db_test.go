package repo_test

import (
	"testing"

	"gitlab.com/ron-liu/cypherscan-server/internal/repo"

	"github.com/cypherium/CypherTestNet/go-cypherium/common"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestTransactionHashIsSearchable(t *testing.T) {
	testOnAnCleanDb(func(db *gorm.DB) {
		tx := repo.Transaction{Hash: repo.Hash(common.Hash{1, 2})}
		db.Debug().Save(&tx)
		var retTxs []repo.Transaction
		db.Debug().Where("hash = ?", repo.Hash(common.Hash{1, 2})).Find(&retTxs)
		assert.Len(t, retTxs, 1)
	})
}
