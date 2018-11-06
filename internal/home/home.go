package home

import (
  "fmt"
  "gitlab.com/ron-liu/cypherscan-server/internal/publisher"
  "math/big"
  "net/http"
  "time"
)

type _TxBlock struct {
  Number    *big.Int
  Txn       int
  CreatedAt time.Time
}

type _KeyBlock struct {
  Number    *big.Int
  CreatedAt time.Time
}
type _MetricValue struct {
  unit   string
  value  float32
  digits int
}
type _Tx struct {
  createdAt time.Time
  value     _MetricValue
  hash      string
  from      string
  to        string
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
  fmt.Fprintf(w, "home")
}
