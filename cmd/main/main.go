package main

import (
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  log "github.com/sirupsen/logrus"
  "gitlab.com/ron-liu/cypherscan-server/internal/blockchain"
  "gitlab.com/ron-liu/cypherscan-server/internal/env"
  "gitlab.com/ron-liu/cypherscan-server/internal/home"
  "gitlab.com/ron-liu/cypherscan-server/internal/publisher"
  "gitlab.com/ron-liu/cypherscan-server/internal/txblock"
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
)

func main() {
  log.SetFormatter(&log.JSONFormatter{})
  log.Info("Evironments:", env.Env)

  util.OpenDb()
  initDb()
  defer util.CloseDb()

  go blockchain.SubscribeNewBlock([]blockchain.BlockHandlers{txblock.SaveBlock, boardcastNewBlock})
  go publisher.StartHub()

  util.CreateRouter(func(r *mux.Router) {
    r.HandleFunc("/home", home.GetHome).Methods("GET")
    r.HandleFunc("/ws", home.HanderForBrowser)
  })
}

func initDb() {
  util.Run(func(db *gorm.DB) error {
    db.AutoMigrate(&txblock.TxBlock{}, &txblock.Transaction{})
    return nil
  })
}

func boardcastNewBlock(block *types.Block) error {
  publisher.Broadcast(home.Payload{
    TxBlocks: []home.TxBlock{*home.TransformTxBlockToFrontend(block)},
    Txs: func() []home.Tx {
      ret := make([]home.Tx, 0, home.TransactionCount)
      for _, transaction := range block.Transactions()[max(0, len(block.Transactions())-home.TransactionCount):] {
        ret = append(ret, *home.TransformTxToFrontend(transaction, block))
      }
      return ret
    }(),
    KeyBlocks: []home.KeyBlock{},
    Metrics:   []home.Metric{},
  })
  return nil
}

func max(a, b int) int {
  if a > b {
    return a
  }
  return b
}
