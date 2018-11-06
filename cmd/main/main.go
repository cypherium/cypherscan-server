package main

import (
  "fmt"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "gitlab.com/ron-liu/cypherscan-server/internal/env"
  "gitlab.com/ron-liu/cypherscan-server/internal/home"
  "gitlab.com/ron-liu/cypherscan-server/internal/publisher"
  "gitlab.com/ron-liu/cypherscan-server/internal/txblock"
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
  "log"
  "net/http"
)

func initDb() {
  util.Run(func(db *gorm.DB) error {
    db.AutoMigrate(&txblock.TxBlock{}, &txblock.Transaction{})
    return nil
  })
}

func main() {
  fmt.Println("Evironments:", env.Env)
  util.OpenDb()
  initDb()
  defer util.CloseDb()
  home.SubscribeNewBlock()
  publisher.StartHub()
  router := mux.NewRouter()
  router.HandleFunc("/home", home.GetHome).Methods("GET")
  router.HandleFunc("/ws", home.HanderForBrowser)
  log.Fatal(http.ListenAndServe(":8844", router))

}
