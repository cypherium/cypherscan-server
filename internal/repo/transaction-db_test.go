package repo_test

import (
	"testing"

	"gitlab.com/ron-liu/cypherscan-server/internal/repo"

	"github.com/cypherium/CypherTestNet/go-cypherium/common"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestTransactionHashFromToIsSearchable(t *testing.T) {
	testOnAnCleanDb(func(db *gorm.DB) {
		hash := repo.Hash(common.Hash{1, 2})
		from := repo.Address(common.Address{3, 4})
		to := repo.Address(common.Address{5, 6})
		nonExistedHash := repo.Hash(common.Hash{})
		nonExistedAddress := repo.Address(common.Address{})

		tx := repo.Transaction{
			Hash: hash,
			From: from,
			To:   to,
		}
		db.Debug().Save(&tx)
		var retTxs []repo.Transaction

		db.Debug().Where("hash = ?", hash).Find(&retTxs)
		assert.Len(t, retTxs, 1)

		db.Debug().Where("\"from\" = ?", from).Find(&retTxs)
		assert.Len(t, retTxs, 1)

		db.Debug().Where("\"to\" = ?", to).Find(&retTxs)
		assert.Len(t, retTxs, 1)

		db.Debug().Where("hash = ?", nonExistedHash).Find(&retTxs)
		assert.Len(t, retTxs, 0)

		db.Debug().Where("\"from\" = ?", nonExistedAddress).Find(&retTxs)
		assert.Len(t, retTxs, 0)

		db.Debug().Where("\"to\" = ?", nonExistedHash).Find(&retTxs)
		assert.Len(t, retTxs, 0)
	})
}
