package home

import (
  "encoding/json"
  "fmt"
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/jinzhu/gorm"
  // log "github.com/sirupsen/logrus"
  "gitlab.com/ron-liu/cypherscan-server/internal/publisher"
  "gitlab.com/ron-liu/cypherscan-server/internal/txblock"
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
  "net/http"
  "time"
)

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
  Number    txblock.UInt64 `json:"number"`
  Txn       int            `json:"txn"`
  CreatedAt time.Time      `json:"createdAt"`
}

// KeyBlock is the key block type transfore to frontend in home page
type KeyBlock struct {
  Number    txblock.UInt64
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
  CreatedAt time.Time      `json:"createdAt"`
  Value     txblock.BigInt `json:"value"`
  Hash      string         `json:"hash"`
  From      string         `json:"from"`
  To        string         `json:"to"`
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

// HanderForBrowser is Websocket handler for browser
func HanderForBrowser(w http.ResponseWriter, r *http.Request) {
  publisher.ServeWebsocket(w, r)
}

// GetHome is to get the initial home contents
func GetHome(w http.ResponseWriter, r *http.Request) {
  fmt.Println("starting getting home")
  var txBlocks []txblock.TxBlock
  var transactions []txblock.Transaction
  util.Run(func(db *gorm.DB) error {
    db.Select([]string{"number", "txn", "time"}).Order("time desc").Limit(TxBlockCount).Find(&txBlocks)
    db.Preload("Block", func(db *gorm.DB) *gorm.DB {
      return db.Select([]string{"time", "hash"})
    }).Select([]string{"block_hash", "value", "hash", "\"from\"", "\"to\""}).Order("transaction_index desc").Limit(TransactionCount).Find(&transactions)
    return nil
  })

  payload := Payload{
    Metrics: []Metric{},
    TxBlocks: func() []TxBlock {
      ret := make([]TxBlock, 0, len(txBlocks))
      for _, b := range txBlocks {
        ret = append(ret, TxBlock{b.Number, b.Txn, b.Time})
      }
      return ret
    }(),
    KeyBlocks: []KeyBlock{},
    Txs: func() []Tx {
      ret := make([]Tx, 0, len(transactions))
      for _, t := range transactions {
        ret = append(ret, Tx{
          t.Block.Time,
          t.Value,
          t.Hash.Hex(),
          t.From.Hex(),
          t.To.Hex(),
        })
      }
      return ret
    }(),
  }
  response, err := json.Marshal(payload)
  if err != nil {
    fmt.Println("error occurs")
  }
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, string(response))
}

// TransformTxBlockToFrontend is to transform Block to frontend json
func TransformTxBlockToFrontend(block *types.Block) *TxBlock {
  return &TxBlock{
    Number:    txblock.UInt64(block.Number().Uint64()),
    Txn:       len(block.Transactions()),
    CreatedAt: time.Unix(block.Time().Int64(), 0),
  }
}

// TransformTxToFrontend is to transform Tx to frontend json
func TransformTxToFrontend(tx *types.Transaction, block *types.Block) *Tx {
  return &Tx{
    CreatedAt: time.Unix(block.Time().Int64(), 0),
    Value:     txblock.BigInt(*tx.Value()),
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
