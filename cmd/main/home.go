package main

import (

	// "encoding/json"

	"net/http"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	"github.com/cypherium/CypherTestNet/go-cypherium/crypto"

	// log "github.com/sirupsen/logrus"

	"gitlab.com/ron-liu/cypherscan-server/internal/repo"

	"time"
)

const (
	BlocksPageSize    = 3
	KeyBlocksPageSize = 3
	TxsPageSize       = 5
)

func getHome(a *App, w http.ResponseWriter, r *http.Request) {
	blockLatestNumber, err := a.blocksFetcher.GetLatestBlockNumber()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	txBlocks, err := a.repo.GetBlocks(&repo.BlockSearchContdition{Scenario: repo.HomePage, StartWith: blockLatestNumber, PageSize: BlocksPageSize})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	keyBlockLatestNumber, err := a.blocksFetcher.GetLatestKeyBlockNumber()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	keyBlocks, err := a.repo.GetKeyBlocks(&repo.BlockSearchContdition{Scenario: repo.HomePage, PageSize: KeyBlocksPageSize, StartWith: keyBlockLatestNumber})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	transactions, err := a.repo.GetTransactions(&repo.TransactionSearchCondition{Scenario: repo.HomePage, PageSize: TxsPageSize})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	payload := HomePayload{
		Metrics: []HomeMetric{},
		TxBlocks: func() []HomeTxBlock {
			ret := make([]HomeTxBlock, 0, len(txBlocks))
			for _, b := range txBlocks {
				ret = append(ret, HomeTxBlock{b.Number, b.Txn, b.Time})
			}
			return ret
		}(),
		KeyBlocks: func() []HomeKeyBlock {
			ret := make([]HomeKeyBlock, 0, len(keyBlocks))
			for _, b := range keyBlocks {
				ret = append(ret, HomeKeyBlock{b.Number, b.Time})
			}
			return ret
		}(),
		Txs: func() []HomeTx {
			ret := make([]HomeTx, 0, len(transactions))
			for _, t := range transactions {
				ret = append(ret, HomeTx{
					t.Block.Time,
					t.Value,
					t.Hash.Hex(),
					t.From,
					t.To,
				})
			}
			return ret
		}(),
	}
	respondWithJSON(w, http.StatusOK, payload)
}

func getHomeTxsFromBlock(block *types.Block, count int) []HomeTx {
	ret := make([]HomeTx, 0, count)
	for _, transaction := range block.Transactions()[max(0, len(block.Transactions())-count):] {
		ret = append(ret, *transformTxToFrontend(transaction, block))
	}
	return ret
}

func transformTxBlocksToFrontendMessage(blocks []*types.Block) *HomePayload {
	txBlocks := make([]HomeTxBlock, 0, len(blocks))
	for _, b := range blocks {
		txBlocks = append(txBlocks, *transformTxBlockToFrontend(b))
	}
	txs := make([]HomeTx, 0, TransactionCount)
	for i := len(blocks) - 1; i >= 0; i-- {
		currentBlock := blocks[i]
		currentTxs := getHomeTxsFromBlock(currentBlock, TransactionCount-len(txs))
		txs = append(txs, currentTxs...)
		if len(txs) >= TransactionCount {
			break
		}
	}
	return &HomePayload{
		TxBlocks:  txBlocks,
		Txs:       txs,
		KeyBlocks: []HomeKeyBlock{},
		Metrics:   []HomeMetric{},
	}
}

func transformKeyBlockToFrontendMessage(block *types.KeyBlockHeader) *HomePayload {
	return &HomePayload{
		TxBlocks:  []HomeTxBlock{},
		Txs:       []HomeTx{},
		Metrics:   []HomeMetric{},
		KeyBlocks: []HomeKeyBlock{*transformKeyBlockToFrontend(block)},
	}
}

const (
	//TxBlockCount is total block number need to return
	TxBlockCount = 3
	//KeyBlockCount is total block number need to return
	KeyBlockCount = 3
	//TransactionCount is total block number need to return
	TransactionCount = 5
)

// HomeTxBlock is the type transfor to frontend in home page
type HomeTxBlock struct {
	Number    int64     `json:"number"`
	Txn       int       `json:"txn"`
	CreatedAt time.Time `json:"createdAt"`
}

// HomeKeyBlock is the key block type transfore to frontend in home page
type HomeKeyBlock struct {
	Number    int64     `json:"number"`
	CreatedAt time.Time `json:"createdAt"`
}

// MetricValue is the MetricValue type transfore to frontend in home page
type MetricValue struct {
	unit   string
	value  float32
	digits int
}

// HomeTx is the HomeTx type trransfore to frontend in home page
type HomeTx struct {
	CreatedAt time.Time    `json:"createdAt"`
	Value     repo.UInt64  `json:"value"`
	Hash      string       `json:"hash"`
	From      repo.Address `json:"from"`
	To        repo.Address `json:"to"`
}

// HomeMetric is the HomeMetric type transfore to frontend in home page
type HomeMetric struct {
	key       string
	name      string
	value     MetricValue
	needGraph bool
}

// HomePayload is the HomePayload type transfore to fronent in home page
type HomePayload struct {
	Metrics   []HomeMetric   `json:"metrics"`
	TxBlocks  []HomeTxBlock  `json:"txBlocks"`
	KeyBlocks []HomeKeyBlock `json:"keyBlocks"`
	Txs       []HomeTx       `json:"txs"`
}

func transformTxBlockToFrontend(block *types.Block) *HomeTxBlock {
	return &HomeTxBlock{
		Number:    block.Number().Int64(),
		Txn:       len(block.Transactions()),
		CreatedAt: time.Unix(0, block.Time().Int64()),
	}
}

func transformKeyBlockToFrontend(block *types.KeyBlockHeader) *HomeKeyBlock {
	return &HomeKeyBlock{
		Number:    block.Number.Int64(),
		CreatedAt: time.Unix(block.Time.Int64(), 0),
	}
}

func transformTxToFrontend(tx *types.Transaction, block *types.Block) *HomeTx {
	return &HomeTx{
		CreatedAt: time.Unix(0, block.Time().Int64()),
		Value:     repo.UInt64(tx.Value().Uint64()),
		Hash:      tx.Hash().Hex(),
		From:      repo.Address(crypto.PubKeyToAddressCypherium(tx.SenderKey())),
		To:        repo.Address(*tx.To()),
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
