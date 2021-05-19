package repo

import (
	"errors"
	"fmt"
	"github.com/cypherium/cypherBFT-P/go-cypherium/core/types"
	"github.com/cypherium/cypherscan-server/internal/util"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"math"
	"reflect"
)

// Get is the interface to get saved information
type Get interface {
	GetBlocks(condition *BlockSearchContdition) ([]TxBlock, error)
	GetBlock(number int64) (*TxBlock, error)
	GetBlockByHash(hash Hash) (*TxBlock, error)
	GetKeyBlock(number int64) (*KeyBlock, error)
	GetKeyBlockByHash(hash Hash) (*KeyBlock, error)
	GetKeyBlocks(condition *BlockSearchContdition) ([]KeyBlock, error)
	GetTransactions(condition *TransactionSearchCondition) ([]Transaction, error)
	GetTransaction(hash Hash) (*Transaction, error)
	GetLocalHighestKeyBlock() (*KeyBlock, error)
	GetLocalHighestBlock() (*TxBlock, error)
	QueryAddress(request *QueryAddressRequest) (*QueryResult, error)
}

// BlockSaver is the interface contains SaveBlock
type BlockSaver interface {
	SaveBlock(block *types.Block) error
	SaveKeyBlock(block *types.KeyBlock) error
	GetLocalHighestKeyBlock() (*KeyBlock, error)
	GetLocalHighestBlock() (*TxBlock, error)
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
		db.Model(&TxBlock{}).AddIndex("idx_block_number", "number")
		db.Model(&Transaction{}).AddIndex("idx_tx_hash", "hash")
		db.Model(&Transaction{}).AddIndex("idx_tx_block_number", "block_number")
		db.Model(&KeyBlock{}).AddIndex("idx_key_block_number", "number")

		return nil
	})
}

// SaveBlock is to save blocks into db
func (repo *Repo) SaveBlock(block *types.Block) error {
	record := transformBlockToDbRecord(block)
	if block.Number().Int64() > 1 {
		getBlock, _ := repo.GetBlock(record.Number)
		if reflect.DeepEqual(getBlock, record) {
			return errors.New("txBlock exist")
		}
	}
	repo.dbRunner.Run(func(db *gorm.DB) error {
		db.Create(record)
		return nil
	})
	log.Infof("SaveBlock number %d", block.Number())
	return nil
}

// SaveKeyBlock is to save key block into db
func (repo *Repo) SaveKeyBlock(block *types.KeyBlock) error {

	record := transferKeyBlockHeaderToDbRecord(block)
	if block.Number().Int64() > 1 {
		getBlock, _ := repo.GetKeyBlock(record.Number)
		if reflect.DeepEqual(getBlock, record) {
			return errors.New("keyBlock exist")
		}
	}
	repo.dbRunner.Run(func(db *gorm.DB) error {
		db.Create(record)
		return nil
	})
	log.Infof("SaveKeyBlock number %d", block.Number())
	return nil
}

func (repo *Repo) GetLocalHighestKeyBlock() (*KeyBlock, error) {

	var keyBlock KeyBlock
	err := repo.dbRunner.Run(func(db *gorm.DB) error {
		db.Last(&keyBlock)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &keyBlock, nil
}

func (repo *Repo) GetLocalHighestBlock() (*TxBlock, error) {

	var txBlock TxBlock
	err := repo.dbRunner.Run(func(db *gorm.DB) error {
		db.Last(&txBlock)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &txBlock, nil
}

// GetBlocks is
func (repo *Repo) GetBlocks(condition *BlockSearchContdition) ([]TxBlock, error) {
	var txBlocks []TxBlock
	log.Infof("GetBlocks  %d", condition)
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
	log.Infof("GetBlock by number %d", number)
	err := repo.dbRunner.Run(func(db *gorm.DB) error {
		return db.Where(whereStatment, whereArgs).Order("time desc").Limit(1).Find(&txBlocks).Error
	})
	if err != nil {
		return nil, err
	}
	if len(txBlocks) <= 0 {
		return nil, &util.MyError{Message: fmt.Sprintf("No Block(number=%d) found in Db", number)}
	}
	return &txBlocks[0], nil
}

func (repo *Repo) GetBlockByHash(hash Hash) (*TxBlock, error) {
	var txBlocks []TxBlock
	whereStatment, whereArgs := func() (string, []interface{}) {
		return "hash = ?", []interface{}{hash}
	}()
	log.Infof("GetBlock by hash %d", hash)
	err := repo.dbRunner.Run(func(db *gorm.DB) error {
		return db.Where(whereStatment, whereArgs).Order("time desc").Limit(1).Find(&txBlocks).Error
	})
	if err != nil {
		return nil, err
	}
	if len(txBlocks) <= 0 {
		return nil, &util.MyError{Message: fmt.Sprintf("No Block(hash=%d) found in Db", hash)}
	}
	return &txBlocks[0], nil
}

// GetKeyBlocks is
func (repo *Repo) GetKeyBlocks(condition *BlockSearchContdition) ([]KeyBlock, error) {
	var keyBlocks []KeyBlock
	log.Infof("GetKeyBlocks  %d", condition)
	pageSize := getPageSizeDefault(condition.PageSize)
	columns := getColumnsByScenario(keyBlockColumnsConfig, condition.Scenario)
	whereStatment, whereArgs := getWhere(condition.StartWith, pageSize)
	return keyBlocks, repo.dbRunner.Run(func(db *gorm.DB) error {
		return db.Where(whereStatment, whereArgs...).Select(columns).Order("time desc").Limit(pageSize).Find(&keyBlocks).Error
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
	log.Infof("GetKeyBlock  %d", number)
	err := repo.dbRunner.Run(func(db *gorm.DB) error {
		return db.Where(whereStatment, whereArgs).Order("time desc").Limit(1).Find(&keyBlocks).Error
	})
	if err != nil {
		return nil, err
	}
	if len(keyBlocks) <= 0 {
		return nil, &util.MyError{Message: fmt.Sprintf("No Block(number=%d) found in Db", number)}
	}
	return &keyBlocks[0], nil
}

func (repo *Repo) GetKeyBlockByHash(hash Hash) (*KeyBlock, error) {
	var keyBlocks []KeyBlock
	whereStatment, whereArgs := func() (string, []interface{}) {
		return "hash = ?", []interface{}{hash}
	}()
	log.Infof("GetKeyBlockByHash  %d", hash)
	err := repo.dbRunner.Run(func(db *gorm.DB) error {
		return db.Where(whereStatment, whereArgs).Order("time desc").Limit(1).Find(&keyBlocks).Error
	})
	if err != nil {
		return nil, err
	}
	if len(keyBlocks) <= 0 {
		return nil, &util.MyError{Message: fmt.Sprintf("No Block(hash=%d) found in Db", hash)}
	}
	return &keyBlocks[0], nil
}

func (repo *Repo) GetTransactions(condition *TransactionSearchCondition) ([]Transaction, error) {
	var txs []Transaction

	pageSize := getPageSizeDefault(condition.PageSize)
	columns := getColumnsByScenario(transactionColumnsConfig, condition.Scenario)
	skip := condition.Skip
	log.Info("GetTransactions BlockNumber", condition.BlockNumber)
	whereStatment, whereArgs := func() (string, []interface{}) {
		if condition.BlockNumber <= 0 {
			return "block_number >= 0", []interface{}{}
		}
		return "block_number = ?", []interface{}{condition.BlockNumber}
	}()
	err := repo.dbRunner.Run(func(db *gorm.DB) error {
		return db.Preload("Block", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"time", "number"})
		}).Where(whereStatment, whereArgs...).Select(columns).Order("id desc").Offset(skip).Limit(pageSize).Find(&txs).Error
	})
	if err != nil {
		return nil, err
	}
	if len(txs) <= 0 {
		return nil, &util.MyError{Message: fmt.Sprintf("No Tranasction(number=%v) found in Db", condition.BlockNumber)}
	}
	var preTransaction Transaction
	var tempTransaction []Transaction
	for _, t := range txs {
		if !reflect.DeepEqual(t, preTransaction) {
			preTransaction = t
			tempTransaction = append(tempTransaction, t)
		}
	}
	txs = tempTransaction
	return txs, nil
}

// GetTransaction is
func (repo *Repo) GetTransaction(hash Hash) (*Transaction, error) {
	var txs []Transaction
	//xlog.Info("GetTransaction")
	log.Infof("GetTransaction  %d", hash)
	whereStatment, whereArgs := func() (string, []interface{}) {
		return "hash = ?", []interface{}{hash}
	}()
	err := repo.dbRunner.Run(func(db *gorm.DB) error {
		return db.Preload("Block", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"time", "number"})
		}).Where(whereStatment, whereArgs).Find(&txs).Error
	})
	if err != nil {
		return nil, err
	}
	if len(txs) <= 0 {
		return nil, &util.MyError{Message: fmt.Sprintf("No Tranasction(number=%v) found in Db", hash)}
	}
	//log.Info("GetTransaction  BlockNumber",&txs[0].BlockNumber)
	return &txs[0], nil
}

var blockColumnsConfig = map[Scenario][]string{
	HomePage: []string{"number", "txn", "time"},
	ListPage: []string{"number", "txn", "time", "gas_used", "gas_limit", "key_signature"},
}
var keyBlockColumnsConfig = map[Scenario][]string{
	HomePage: []string{"number", "time"},
	ListPage: []string{"number", "time", "difficulty"},
}
var transactionColumnsConfig = map[Scenario][]string{
	HomePage: []string{"value", "hash", "\"from\"", "\"to\"", "block_number"},
	ListPage: []string{"value", "hash", "\"from\"", "\"to\"", "block_number"},
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
	var start int64
	if (startWith - pageSize + 1) >= 0 {
		start = startWith - pageSize + 1
	} else {
		start = 0
	}
	return "number BETWEEN ? AND ?", []interface{}{start, startWith}
}
