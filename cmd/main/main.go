package main

import (
  "fmt"
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
  "net/http"
)

func initDb() {
  util.Run(func(db *gorm.DB) error {
    db.AutoMigrate(&txblock.TxBlock{}, &txblock.Transaction{})
    return nil
  })
}

func main() {
  log.SetFormatter(&log.JSONFormatter{})
  fmt.Println("Evironments:", env.Env)
  util.OpenDb()
  initDb()
  defer util.CloseDb()
  blockchain.SubscribeNewBlock()
  publisher.StartHub()
  router := mux.NewRouter()
  router.HandleFunc("/home", home.GetHome).Methods("GET")
  router.HandleFunc("/ws", home.HanderForBrowser)
  log.Fatal(http.ListenAndServe(":8844", router))

}
