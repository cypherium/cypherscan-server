package main

import (
	// "encoding/json"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"

	// log "github.com/sirupsen/logrus"

	"gitlab.com/ron-liu/cypherscan-server/internal/repo"

	"time"
)

func transformTxBlockToFrontendMessage(block *types.Block) *HomePayload {
	return &HomePayload{
		TxBlocks: []HomeTxBlock{*transformTxBlockToFrontend(block)},
		Txs: func() []HomeTx {
			ret := make([]HomeTx, 0, TransactionCount)
			for _, transaction := range block.Transactions()[max(0, len(block.Transactions())-TransactionCount):] {
				ret = append(ret, *transformTxToFrontend(transaction, block))
			}
			return ret
		}(),
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
	CreatedAt time.Time   `json:"createdAt"`
	Value     repo.BigInt `json:"value"`
	Hash      string      `json:"hash"`
	From      string      `json:"from"`
	To        string      `json:"to"`
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
		CreatedAt: time.Unix(block.Time().Int64(), 0),
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
		CreatedAt: time.Unix(block.Time().Int64(), 0),
		Value:     repo.BigInt(*tx.Value()),
		Hash:      tx.Hash().Hex(),
		From:      "",
		To:        "",
		// To: func() string {
		//   to := tx.To()
		//   if tx == nil {
		//     return ""
		//   }
		//   log.Infoln("to.Hex()", to)
		//   return to.Hex()
		// }(),
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
