package repo_test

import (
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/stretchr/testify/assert"

	"github.com/cypherium/cypherscan-server/internal/repo"
)

func TestSaveBlockWithTransaction(t *testing.T) {
	keySignature := repo.Bytes([]byte{1, 2})
	nonExistedKeySignature := repo.Bytes([]byte{1, 2, 3})
	testOnAnCleanDb(func(db *gorm.DB) {
		block := repo.TxBlock{Number: 1, Transactions: []repo.Transaction{repo.Transaction{}}, KeySignature: keySignature}
		db.Debug().Save(&block)
		var retBlocks []repo.TxBlock
		var retTransactions []repo.Transaction
		db.Debug().Find(&retBlocks)
		db.Debug().Find(&retTransactions)
		assert.Len(t, retBlocks, 1)
		assert.Equal(t, keySignature, retBlocks[0].KeySignature)
		assert.Len(t, retTransactions, 1)

		db.Debug().Where("key_signature = ?", keySignature).Find(&retBlocks)
		assert.Len(t, retBlocks, 1)

		db.Debug().Where("key_signature = ?", nonExistedKeySignature).Find(&retBlocks)
		assert.Len(t, retBlocks, 0)
	})
}
