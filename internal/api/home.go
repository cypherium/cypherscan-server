package api

import (
	// "encoding/json"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"

	// log "github.com/sirupsen/logrus"

	"gitlab.com/ron-liu/cypherscan-server/internal/repo"

	"time"
)

// TransformTxBlockToFrontendMessage is to transform eth block type to message will broadcast to browsers
func TransformTxBlockToFrontendMessage(block *types.Block) *Payload {
	return &Payload{
		TxBlocks: []TxBlock{*transformTxBlockToFrontend(block)},
		Txs: func() []Tx {
			ret := make([]Tx, 0, TransactionCount)
			for _, transaction := range block.Transactions()[max(0, len(block.Transactions())-TransactionCount):] {
				ret = append(ret, *transformTxToFrontend(transaction, block))
			}
			return ret
		}(),
		KeyBlocks: []KeyBlock{},
		Metrics:   []Metric{},
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

// TxBlock is the type transfor to frontend in home page
type TxBlock struct {
	Number    int64     `json:"number"`
	Txn       int       `json:"txn"`
	CreatedAt time.Time `json:"createdAt"`
}

// KeyBlock is the key block type transfore to frontend in home page
type KeyBlock struct {
	Number    int64
	CreatedAt time.Time
}

// MetricValue is the MetricValue type transfore to frontend in home page
type MetricValue struct {
	unit   string
	value  float32
	digits int
}

// Tx is the Tx type trransfore to frontend in home page
type Tx struct {
	CreatedAt time.Time   `json:"createdAt"`
	Value     repo.BigInt `json:"value"`
	Hash      string      `json:"hash"`
	From      string      `json:"from"`
	To        string      `json:"to"`
}

// Metric is the Metric type transfore to frontend in home page
type Metric struct {
	key       string
	name      string
	value     MetricValue
	needGraph bool
}

// Payload is the Payload type transfore to fronent in home page
type Payload struct {
	Metrics   []Metric   `json:"metrics"`
	TxBlocks  []TxBlock  `json:"txBlocks"`
	KeyBlocks []KeyBlock `json:"keyBlocks"`
	Txs       []Tx       `json:"txs"`
}

func transformTxBlockToFrontend(block *types.Block) *TxBlock {
	return &TxBlock{
		Number:    block.Number().Int64(),
		Txn:       len(block.Transactions()),
		CreatedAt: time.Unix(block.Time().Int64(), 0),
	}
}

func transformTxToFrontend(tx *types.Transaction, block *types.Block) *Tx {
	return &Tx{
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
