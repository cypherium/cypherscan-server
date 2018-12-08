package blockchain

import (
  "context"
  "github.com/ethereum/go-ethereum/ethclient"
  log "github.com/sirupsen/logrus"
  "math/big"

  "github.com/ethereum/go-ethereum/core/types"
)

const clientPoolSize = 10

// GetBlocksByNumbers will accept array of numbers and retrieve the blocks concurrently
func GetBlocksByNumbers(numbers []*big.Int) []*types.Block {
  gather := make(chan (*types.Block), clientPoolSize)
  for _, number := range numbers {
    go Run(func(cli *ethclient.Client) error {
      block, err := cli.BlockByNumber(context.Background(), number)
      if err == nil {
        log.Error("Cannot get block by hash", number, err)
        gather <- nil
      } else {
        gather <- block
      }
      return nil
    })
  }
  out := make([]*types.Block, len(numbers))
  for range numbers {
    out = append(out, <-gather)
  }
  return out
}
