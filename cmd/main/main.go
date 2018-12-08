package main

import (
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  log "github.com/sirupsen/logrus"
  "gitlab.com/ron-liu/cypherscan-server/internal/api"
  "gitlab.com/ron-liu/cypherscan-server/internal/blockchain"
  "gitlab.com/ron-liu/cypherscan-server/internal/env"
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

  blockchain.Connect()
  go blockchain.SubscribeNewBlock([]blockchain.BlockHandlers{txblock.SaveBlock, boardcastNewBlock})
  go publisher.StartHub()

  util.CreateRouter(func(r *mux.Router) {
    r.HandleFunc("/home", api.GetHome).Methods("GET")
    r.HandleFunc("/ws", api.HanderForBrowser)
  })
}

func initDb() {
  util.Run(func(db *gorm.DB) error {
    db.AutoMigrate(&txblock.TxBlock{}, &txblock.Transaction{})
    return nil
  })
}

func boardcastNewBlock(block *types.Block) error {
  publisher.Broadcast(api.TransformTxBlockToFrontendMessage(block))
  return nil
}
