package main

import (
  "context"
  "fmt"
  "log"
  "math/big"
  "time"

  "github.com/ethereum/go-ethereum/rpc"
  "github.com/gin-gonic/gin"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "gitlab.com/ron-liu/cypherscan-server/internal/env"
  "gitlab.com/ron-liu/cypherscan-server/internal/home"
  "gitlab.com/ron-liu/cypherscan-server/internal/txblock"
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
)

func initDb() {
  db := util.OpenDb()
  db.AutoMigrate(&txblock.TxBlock{})
  defer db.Close()
}

type Block struct {
  Number string
}

func main() {
  client, err := rpc.Dial("wss://mainnet.infura.io/ws")
  if err != nil {
    fmt.Println(err)
    return
  }
  subch := make(chan Block)

  // Ensure that subch receives the latest block.
  go func() {
    for i := 0; ; i++ {
      if i > 0 {
        time.Sleep(2 * time.Second)
      }
      subscribeBlocks(client, subch)
    }
  }()

  // Print events from the subscription as they arrive.
  for block := range subch {
    x := new(big.Int)
    x, ok := x.SetString(block.Number[2:], 16)
    if !ok {
      log.Fatal("error: ", block.Number)
    }
    fmt.Println("latest block:", block.Number, x)
    // fmt.Println("latest block:", block.Number, "diff", block.Difficulty, "txHash", block.TxHash, "gas limit", block.GasLimit, "gs used", block.GasUsed, block.Nonce)
  }

  fmt.Println("Evironments:", env.Env)
  routers := gin.Default()
  routers.GET("/home", home.GetHome)
  routers.Run()
}

func subscribeBlocks(client *rpc.Client, subch chan Block) {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  // Subscribe to new blocks.
  sub, err := client.EthSubscribe(ctx, subch, "newHeads")
  if err != nil {
    fmt.Println("subscribe error:", err)
    return
  }

  // The connection is established now.
  // Update the channel with the current block.
  var lastBlock Block
  if err := client.CallContext(ctx, &lastBlock, "eth_getBlockByNumber", "latest", true); err != nil {
    fmt.Println("can't get latest block:", err)
    return
  }
  subch <- lastBlock

  // The subscription will deliver events to the channel. Wait for the
  // subscription to end for any reason, then loop around to re-establish
  // the connection.
  fmt.Println("connection lost: ", <-sub.Err())
}
