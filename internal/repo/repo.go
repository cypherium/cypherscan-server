package repo

import (
	"fmt"
	"math"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	"github.com/jinzhu/gorm"
	"gitlab.com/ron-liu/cypherscan-server/internal/util"
)

// Get is the interface to get saved information
type Get interface {
	GetBlocks(condition *BlockSearchContdition) ([]TxBlock, error)
	GetBlock(number int64) (*TxBlock, error)
	GetKeyBlock(number int64) (*KeyBlock, error)
	GetKeyBlocks(condition *BlockSearchContdition) ([]KeyBlock, error)
	GetTransactions(condition *TransactionSearchCondition) ([]Transaction, error)
}

// BlockSaver is the interface contains SaveBlock
type BlockSaver interface {
	SaveBlock(block *types.Block) error
	SaveKeyBlock(block *types.KeyBlockHeader) error
}

// Repo is the database access layer
type Repo struct {
	dbRunner util.DbRunner
}

// NewRepo is the constructor to create Repo
func NewRepo(dbRunner util.DbRunner) *Repo {
	return &Repo{dbRunner}
}

// InitDb is to create the db structure
func (repo *Repo) InitDb() {
	repo.dbRunner.Run(func(db *gorm.DB) error {
		db.AutoMigrate(&TxBlock{}, &Transaction{}, &KeyBlock{})
		return nil
	})
}

// SaveBlock is to save blocks into db
func (repo *Repo) SaveBlock(block *types.Block) error {
	record := transformBlockToDbRecord(block)
	// util.PrintStructInJSON(record)
	repo.dbRunner.Run(func(db *gorm.DB) error {
		db.NewRecord(record)
		db.Debug().Create(record)
		return nil
	})
	return nil
}

// SaveKeyBlock is to save key block into db
func (repo *Repo) SaveKeyBlock(block *types.KeyBlockHeader) error {
	record := transferKeyBlockHeaderToDbRecord(block)
	repo.dbRunner.Run(func(db *gorm.DB) error {
		db.NewRecord(record)
		db.Debug().Create(record)
		return nil
	})
	return nil
}

// GetBlocks is
func (repo *Repo) GetBlocks(condition *BlockSearchContdition) ([]TxBlock, error) {
	var txBlocks []TxBlock
	pageSize := getPageSizeDefault(condition.PageSize)
	columns := getColumnsByScenario(blockColumnsConfig, condition.Scenario)
	whereStatment, whereArgs := getWhere(condition.StartWith, pageSize)
	return txBlocks, repo.dbRunner.Run(func(db *gorm.DB) error {
		return db.Debug().Where(whereStatment, whereArgs...).Select(columns).Order("time desc").Limit(pageSize).Find(&txBlocks).Error
	})
}

// GetBlock is to get single block by the number if number >=0, otherwise it will get the latest one
func (repo *Repo) GetBlock(number int64) (*TxBlock, error) {
	var txBlocks []TxBlock
	whereStatment, whereArgs := func() (string, []interface{}) {
		if number < 0 {
			return "1=1", []interface{}{}
		}
		return "number = ?", []interface{}{number}
	}()

	err := repo.dbRunner.Run(func(db *gorm.DB) error {
		return db.Debug().Where(whereStatment, whereArgs).Select(getColumnsByScenario(blockColumnsConfig, ListPage)).Order("time desc").Limit(1).Find(&txBlocks).Error
	})
	if err != nil {
		return nil, err
	}
	if len(txBlocks) <= 0 {
		return nil, &util.MyError{Message: fmt.Sprintf("No Block(number=%d) found in Db", number)}
	}
	return &txBlocks[0], nil
}

// GetKeyBlocks is
func (repo *Repo) GetKeyBlocks(condition *BlockSearchContdition) ([]KeyBlock, error) {
	var keyBlocks []KeyBlock
	pageSize := getPageSizeDefault(condition.PageSize)
	columns := getColumnsByScenario(keyBlockColumnsConfig, condition.Scenario)
	whereStatment, whereArgs := getWhere(condition.StartWith, pageSize)
	return keyBlocks, repo.dbRunner.Run(func(db *gorm.DB) error {
		return db.Debug().Where(whereStatment, whereArgs...).Select(columns).Order("time desc").Limit(pageSize).Find(&keyBlocks).Error
	})
}

// GetKeyBlock is to get single key block by the number if number >=0, otherwise it will get the latest one
func (repo *Repo) GetKeyBlock(number int64) (*KeyBlock, error) {
	var keyBlocks []KeyBlock
	whereStatment, whereArgs := func() (string, []interface{}) {
		if number < 0 {
			return "1=1", []interface{}{}
		}
		return "number = ?", []interface{}{number}
	}()

	err := repo.dbRunner.Run(func(db *gorm.DB) error {
		return db.Debug().Where(whereStatment, whereArgs).Select(getColumnsByScenario(blockColumnsConfig, ListPage)).Order("time desc").Limit(1).Find(&keyBlocks).Error
	})
	if err != nil {
		return nil, err
	}
	if len(keyBlocks) <= 0 {
		return nil, &util.MyError{Message: fmt.Sprintf("No Block(number=%d) found in Db", number)}
	}
	return &keyBlocks[0], nil
}

// GetTransactions is
func (repo *Repo) GetTransactions(condition *TransactionSearchCondition) ([]Transaction, error) {
	var txs []Transaction
	pageSize := getPageSizeDefault(condition.PageSize)
	columns := getColumnsByScenario(transactionColumnsConfig, condition.Scenario)
	skip := condition.Skip
	whereStatment, whereArgs := func() (string, []interface{}) {
		if condition.BlockNumber == 0 {
			return "block_number >= 0", []interface{}{}
		}
		return "number block_number = ?", []interface{}{condition.BlockNumber}
	}()
	return txs, repo.dbRunner.Run(func(db *gorm.DB) error {
		return db.Debug().Preload("Block", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"time", "hash"})
		}).Select(whereStatment, whereArgs...).Select(columns).Order("block_number desc, transaction_index desc").Offset(skip).Limit(pageSize).Find(&txs).Error
	})
}

var blockColumnsConfig = map[Scenario][]string{
	HomePage: []string{"number", "txn", "time"},
	ListPage: []string{"number", "txn", "time"},
}
var keyBlockColumnsConfig = map[Scenario][]string{
	HomePage: []string{"number", "time"},
	ListPage: []string{"number", "time"},
}
var transactionColumnsConfig = map[Scenario][]string{
	HomePage: []string{"block_hash", "value", "hash", "\"from\"", "\"to\""},
	ListPage: []string{"block_hash", "value", "hash", "\"from\"", "\"to\""},
}

func getColumnsByScenario(config map[Scenario][]string, scenario Scenario) []string {
	return config[scenario]
}

// Scenario defines the scenario using columns
type Scenario int

const (
	// HomePage is used in homepage and just few columns needed
	HomePage Scenario = 0
	//ListPage is used in list page
	ListPage Scenario = 1
)

const defaultPageSize = 3

// BlockSearchContdition contains search conditions
type BlockSearchContdition struct {
	Scenario  Scenario
	StartWith int64
	PageSize  int
}

// TransactionSearchCondition contains search conditions for transactions
type TransactionSearchCondition struct {
	Scenario    Scenario
	BlockNumber int64
	Skip        int64
	PageSize    int
}

func getPageSizeDefault(pageSize int) int64 {
	if pageSize == 0 {
		return defaultPageSize
	}
	return int64(pageSize)
}

func getWhere(startWith int64, pageSize int64) (string, []interface{}) {
	if startWith < 0 {
		return "number <= ?", []interface{}{math.MaxInt64}
	}
	return "number BETWEEN ? AND ?", []interface{}{startWith - pageSize + 1, startWith}
}
