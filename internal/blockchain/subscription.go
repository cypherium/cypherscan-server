package blockchain

import (
  "context"
  "fmt"
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/ethereum/go-ethereum/ethclient"
  log "github.com/sirupsen/logrus"
  "gitlab.com/ron-liu/cypherscan-server/internal/env"
)

// BlockHandlers is the function type to handle the recieved block
type BlockHandlers func(*types.Block) error

// SubscribeNewBlock is to subscribe new block
func SubscribeNewBlock(handlers []BlockHandlers) {
  fmt.Println("about to connect")
  client, err := ethclient.Dial(env.Env.TsBlockChainWsURL)
  if err != nil {
    log.Fatal("Cannot connect blockchain nodes", err)
    return
  }
  log.Info("Connected to blockchain nodes")
  headers := make(chan *types.Header)

  sub, err := client.SubscribeNewHead(context.Background(), headers)
  if err != nil {
    log.Fatal("Cannot subscribe new heads from blockchain nodes", err)
  }
  log.Info("Subscribed to blockchain nodes")

  for {
    select {
    case err := <-sub.Err():
      log.Fatal(err)
    case header := <-headers:
      block, err := client.BlockByHash(context.Background(), header.Hash())
      if err != nil {
        log.Error("Cannot get block by hash", header.Hash(), err)
        continue
      }
      log.WithFields(log.Fields{
        "Number": block.Number(),
        "txn":    len(block.Transactions()),
      }).Info("A new block generated")
      for _, handler := range handlers {
        handler(block)
      }
    }
  }
}
