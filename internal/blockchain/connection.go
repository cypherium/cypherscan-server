package blockchain

import (
  "github.com/ethereum/go-ethereum/ethclient"
  log "github.com/sirupsen/logrus"
  "gitlab.com/ron-liu/cypherscan-server/internal/env"
)

var client *ethclient.Client

// RunFunc is a type used to call ethclient
type RunFunc func(client *ethclient.Client) error

func getClient() *ethclient.Client {
  return client
}

// Connect to the blockchain and keep the connection in memory
func Connect() {
  _client, err := ethclient.Dial(env.Env.TsBlockChainWsURL)
  if err != nil {
    log.Fatal("Cannot connect blockchain nodes", err)
    return
  }
  log.Info("Connected to blockchain nodes")
  client = _client
}

// Run is to call ethclient with enclosed eth client
func Run(fn RunFunc) error {
  cli := getClient()
  return fn(cli)
}
