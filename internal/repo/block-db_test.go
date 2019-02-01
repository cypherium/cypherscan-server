package repo_test

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/stretchr/testify/assert"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
	"gitlab.com/ron-liu/cypherscan-server/internal/util"
)

func TestSaveBlockWithTransaction(t *testing.T) {
	os.Remove("test.db")
	dbClient, err := util.ConnectDb("sqlite3", nil, nil, "")
	assert.Nil(t, err)
	defer dbClient.Close()
	dbClient.Run(func(db *gorm.DB) error {
		db.AutoMigrate(&repo.TxBlock{}, &repo.Transaction{})
		return nil
	})
	dbClient.Run(func(db *gorm.DB) error {
		block := repo.TxBlock{Number: 1, Transactions: []repo.Transaction{repo.Transaction{}}}
		db.Debug().Save(&block)
		var retBlocks []repo.TxBlock
		var retTransactions []repo.Transaction
		db.Debug().Find(&retBlocks)
		db.Debug().Find(&retTransactions)
		assert.Len(t, retBlocks, 1)
		assert.Len(t, retTransactions, 1)
		return nil
	})

}
