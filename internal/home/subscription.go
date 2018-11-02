package home

import (
  "context"
  "fmt"
  "github.com/ethereum/go-ethereum/rpc"
  "gitlab.com/ron-liu/cypherscan-server/internal/env"
  "gitlab.com/ron-liu/cypherscan-server/internal/txblock"
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
  "math/big"
  "time"
)

type blockInfo struct {
  Number string
}

func subscribeNewBlock() {
  client, err := rpc.Dial(env.Env.TsBlockChainWsURL)
  if err != nil {
    fmt.Println(err)
    return
  }
  subch := make(chan blockInfo)

  // Ensure that subch receives the latest block.
  go func() {
    for {
      subscribeBlocks(client, subch)
      time.Sleep(2 * time.Second)
    }
  }()

  // Print events from the subscription as they arrive.
  for block := range subch {
    fmt.Println("latest block:", block.Number)
    // fmt.Println("latest block:", block.Number, "diff", block.Difficulty, "txHash", block.TxHash, "gas limit", block.GasLimit, "gs used", block.GasUsed, block.Nonce)
  }
}

func transformToTxBlock(b blockInfo) *txblock.TxBlock {
  return &txblock.TxBlock{
    Number: util.Parse(b.Number, new(big.Int)).(*big.Int),
  }
}

func subscribeBlocks(client *rpc.Client, subch chan blockInfo) {
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
  var lastBlock blockInfo
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
