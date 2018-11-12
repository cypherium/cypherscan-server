package home

import (
  "encoding/json"
  "fmt"
  "github.com/jinzhu/gorm"
  log "github.com/sirupsen/logrus"
  "gitlab.com/ron-liu/cypherscan-server/internal/publisher"
  "gitlab.com/ron-liu/cypherscan-server/internal/txblock"
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
  "net/http"
  "time"
)

const (
  //TxBlockCount is total block number need to return
  TxBlockCount = 5
  //KeyBlockCount is total block number need to return
  KeyBlockCount = 5
  //TransactionCount is total block number need to return
  TransactionCount = 5
)

type _TxBlock struct {
  Number    txblock.UInt64 `json:"number"`
  Txn       int            `json:"txn"`
  CreatedAt time.Time      `json:"createdAt"`
}

type _KeyBlock struct {
  Number    txblock.UInt64
  CreatedAt time.Time
}
type _MetricValue struct {
  unit   string
  value  float32
  digits int
}
type _Tx struct {
  CreatedAt time.Time      `json:"createdAt"`
  Value     txblock.BigInt `json:"value"`
  Hash      string         `json:"hash"`
  From      string         `json:"from"`
  To        string         `json:"to"`
}
type _Metric struct {
  key       string
  name      string
  value     _MetricValue
  needGraph bool
}

type home struct {
  Metrics   []_Metric
  TxBlocks  []_TxBlock
  KeyBlocks []_KeyBlock
  Txs       []_Tx
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
    db.Debug().Preload("Block", func(db *gorm.DB) *gorm.DB {
      return db.Select([]string{"time", "hash"})
    }).Select([]string{"block_hash", "value", "hash", "\"from\"", "\"to\""}).Order("transaction_index desc").Limit(TransactionCount).Find(&transactions)
    return nil
  })

  // log.Infof("blocks: %+v\n", txBlocks)
  log.Infof("transactins: %+v\n", transactions)

  payload := home{
    TxBlocks: func() []_TxBlock {
      ret := make([]_TxBlock, 0, len(txBlocks))
      for _, b := range txBlocks {
        ret = append(ret, _TxBlock{b.Number, b.Txn, b.Time})
      }
      return ret
    }(),
    Txs: func() []_Tx {
      ret := make([]_Tx, 0, len(transactions))
      for _, t := range transactions {
        ret = append(ret, _Tx{
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
  // fmt.Println("response", string(response))
}
