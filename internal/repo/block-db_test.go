package repo_test

import (
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/stretchr/testify/assert"

	"github.com/cypherium/cypherscan-server/internal/repo"
)

func TestSaveBlockWithTransaction(t *testing.T) {
	Signature := repo.Bytes([]byte{1, 2})
	nonExistedSignature := repo.Bytes([]byte{1, 2, 3})
	testOnAnCleanDb(func(db *gorm.DB) {
		block := repo.TxBlock{Number: 1, Transactions: []repo.Transaction{repo.Transaction{}}, Signature: Signature}
		db.Debug().Save(&block)
		var retBlocks []repo.TxBlock
		var retTransactions []repo.Transaction
		db.Debug().Find(&retBlocks)
		db.Debug().Find(&retTransactions)
		assert.Len(t, retBlocks, 1)
		assert.Equal(t, Signature, retBlocks[0].Signature)
		assert.Len(t, retTransactions, 1)

		db.Debug().Where("key_signature = ?", Signature).Find(&retBlocks)
		assert.Len(t, retBlocks, 1)

		db.Debug().Where("key_signature = ?", nonExistedSignature).Find(&retBlocks)
		assert.Len(t, retBlocks, 0)
	})
}
