package main

import (
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/gorilla/handlers"
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

func main() {
  log.SetFormatter(&log.JSONFormatter{})
  log.Info("Evironments:", env.Env)

  util.OpenDb()
  initDb()
  defer util.CloseDb()

  go blockchain.SubscribeNewBlock([]blockchain.BlockHandlers{txblock.SaveBlock, boardcastNewBlock})
  go publisher.StartHub()

  createRouter(func(r *mux.Router) {
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

type setupRoute func(*mux.Router)

func createRouter(f setupRoute) {
  r := mux.NewRouter()
  // Handle all preflight request
  r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // fmt.Printf("OPTIONS")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
    w.WriteHeader(http.StatusNoContent)
    return
  })
  r.StrictSlash(true)
  headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
  originsOk := handlers.AllowedOrigins([]string{env.Env.OriginAllowed})
  methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
  f(r)
  log.Fatal(http.ListenAndServe(":8000", handlers.CORS(methodsOk, headersOk, originsOk)(r)))
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
