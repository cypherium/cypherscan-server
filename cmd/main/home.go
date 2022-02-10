package main

import (
	// "encoding/json"

	"github.com/cypherium/cypherBFT/core/types"
	"github.com/cypherium/cypherBFT/crypto"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
	"reflect"

	"github.com/cypherium/cypherscan-server/internal/repo"
	"time"
)

const (
	BlocksPageSize    = 3
	KeyBlocksPageSize = 3
	TxsPageSize       = 5
)

func getHome(a *App, w http.ResponseWriter, r *http.Request) {
	logrus.Info("getHome")
	blockLatestNumber, err := a.blocksFetcher.GetLatestBlockNumber()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var txBlock, preTxBlock *repo.TxBlock
	var txBlocks []repo.TxBlock
	startNumber := blockLatestNumber
	var preTransaction repo.Transaction
	var tempTransaction, transactions []repo.Transaction
	for {
		txBlock, err = a.repo.GetBlock(startNumber)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
		if !reflect.DeepEqual(txBlock, preTxBlock) {
			preTxBlock = txBlock
			txBlocks = append(txBlocks, *txBlock)
		}
		if len(txBlocks) >= BlocksPageSize {
			break
		} else {
			startNumber--
		}
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
	startTransactionNumber := blockLatestNumber
	for {
		transactions, err = a.repo.GetTransactions(&repo.TransactionSearchCondition{BlockNumber: startTransactionNumber, Scenario: repo.HomePage, PageSize: BlocksPageSize})
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
		for _, t := range transactions {
			//logrus.Info("getHome transaction ",t)
			if !reflect.DeepEqual(t, preTransaction) {
				preTransaction = t
				tempTransaction = append(tempTransaction, t)
			}
		}
		if len(tempTransaction) < TxsPageSize {
			startTransactionNumber--
		} else {
			transactions = tempTransaction
			break
		}
	}
	//logrus.Info("getHome transactions len ", len(transactions))
	//var preKeyBlock repo.KeyBlock
	payload := HomePayload{
		Metrics: []HomeMetric{
			HomeMetric{Key: "tps", Name: "TPS", Value: MetricValue{Unit: ""}},
			HomeMetric{Key: "bps", Name: "BPS", Value: MetricValue{Unit: "blocks/sec"}},
			HomeMetric{Key: "key-blocks-nodes", Name: "Validators", Value: MetricValue{Value: 21}},
			HomeMetric{Key: "key-blocks-Diff", Name: "Key Block Diff", Value: MetricValue{}},
			HomeMetric{Key: "tx-blocks-number", Name: "Tx Blocks Number", Value: MetricValue{Value: blockLatestNumber}},
			HomeMetric{Key: "key-blocks-number", Name: "Key Blocks Number", Value: MetricValue{Value: keyBlockLatestNumber}},
		},
		TxBlocks: func() []HomeTxBlock {
			ret := make([]HomeTxBlock, 0, len(txBlocks))
			for _, b := range txBlocks {
				ret = append(ret, HomeTxBlock{b.Number, b.Txn, b.Time})
			}
			return ret
		}(),
		//KeyBlocks: func() []HomeKeyBlock {
		//	ret := make([]HomeKeyBlock, 0, len(keyBlocks))
		//	for _, b := range keyBlocks {
		//		if !reflect.DeepEqual(b, preKeyBlock) {
		//			preKeyBlock = b
		//			ret = append(ret, HomeKeyBlock{b.Number, b.Time})
		//		}
		//	}
		//	return ret
		//}(),
		KeyBlocks: func() []HomeKeyBlock {
			ret := make([]HomeKeyBlock, 0, len(keyBlocks))
			for _, b := range keyBlocks {
				//preKeyBlock = b
				ret = append(ret, HomeKeyBlock{b.Number, b.Time})
			}

			return ret
		}(),
		Txs: func() []HomeTx {
			ret := make([]HomeTx, 0, len(transactions))
			for _, t := range transactions {
				//logrus.Info("getHome transaction ",t)
				if !reflect.DeepEqual(t, preTransaction) {
					preTransaction = t
					ret = append(ret, HomeTx{
						t.Block.Time,
						t.Value,
						t.Hash,
						t.From.String(),
						t.To.String(),
					})
				}
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

type metrics struct {
	currentKeyBlock *types.KeyBlock
}

func transformTxBlocksToFrontendMessage(blocks []*types.Block, metrics metrics) *HomePayload {
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
	totalTxs := int64(0)
	for _, b := range blocks {
		totalTxs += int64(len(b.Transactions()))
	}
	homeMetrics := []HomeMetric{
		HomeMetric{Key: "key-blocks-number", Name: "Key Blocks Number", Value: MetricValue{Value: metrics.currentKeyBlock.Number().Int64()}},
		HomeMetric{Key: "key-blocks-Diff", Name: "Key Block Diff", Value: MetricValue{Unit: "M", Digits: 2, Value: metrics.currentKeyBlock.Difficulty().Int64() / 10000}},
	}
	if len(blocks) > 0 {
		firstBlock := blocks[0]
		lastBlock := blocks[len(blocks)-1]
		tps, bps := func() (int64, int64) {
			ns := (lastBlock.Time().Int64() - firstBlock.Time().Int64())
			if ns == 0 {
				return 0, 0
			}
			return div(totalTxs*int64(math.Pow(10, 9)), ns), div(int64(len(blocks))*int64(math.Pow(10, 9)), ns)
		}()
		homeMetrics = append(
			[]HomeMetric{
				HomeMetric{Key: "tps", Name: "TPS", Value: MetricValue{Value: tps, Unit: ""}},
				HomeMetric{Key: "bps", Name: "BPS", Value: MetricValue{Value: bps, Unit: "blocks/sec"}},
				HomeMetric{Key: "tx-blocks-number", Name: "Tx Blocks Number", Value: MetricValue{Value: lastBlock.Number().Int64()}},
			},
			homeMetrics...,
		)

	}
	// log.Printf("calu tps: total tx: %d, lastBlock time: %v, firstBlock time: %v", totalTxs, lastBlock.Time(), firstBlock.Time())
	return &HomePayload{
		TxBlocks:  txBlocks,
		Txs:       txs,
		KeyBlocks: []HomeKeyBlock{},
		Metrics:   homeMetrics,
	}
}

func div(x, y int64) int64 {
	ret := x / y * 5 / 2
	mod := x % y
	if mod > 1 {
		return ret + 1
	}
	return ret
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
	Unit   string `json:"unit"`
	Value  int64  `json:"value"`
	Digits int    `json:"digits"`
}

// HomeTx is the HomeTx type trransfore to frontend in home page
type HomeTx struct {
	CreatedAt time.Time `json:"createdAt"`
	Value     string    `json:"value"`
	Hash      repo.Hash `json:"hash"`
	From      string    `json:"from"`
	To        string    `json:"to"`
}

// HomeMetric is the HomeMetric type transfore to frontend in home page
type HomeMetric struct {
	Key       string      `json:"key"`
	Name      string      `json:"name"`
	Value     MetricValue `json:"value"`
	NeedGraph bool        `json:"needGraph"`
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
		Value:     tx.Value().String(),
		Hash:      repo.Hash(tx.Hash()),
		From:      repo.Address(crypto.PubKeyToAddressCypherium(tx.SenderKey())).String(),
		To:        repo.Address(*tx.To()).String(),
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
