package repo_test

import (
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/stretchr/testify/assert"

	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
)

func TestSaveBlockWithTransaction(t *testing.T) {
	testOnAnCleanDb(func(db *gorm.DB) {
		block := repo.TxBlock{Number: 1, Transactions: []repo.Transaction{repo.Transaction{}}}
		db.Debug().Save(&block)
		var retBlocks []repo.TxBlock
		var retTransactions []repo.Transaction
		db.Debug().Find(&retBlocks)
		db.Debug().Find(&retTransactions)
		assert.Len(t, retBlocks, 1)
		assert.Len(t, retTransactions, 1)
	})
}
