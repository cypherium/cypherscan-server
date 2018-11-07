package blockchain

import (
  "context"
  "fmt"
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/ethereum/go-ethereum/ethclient"
  "github.com/jinzhu/gorm"
  log "github.com/sirupsen/logrus"
  "gitlab.com/ron-liu/cypherscan-server/internal/env"
  "gitlab.com/ron-liu/cypherscan-server/internal/txblock"
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
  "time"
)

// SubscribeNewBlock is to subscribe new block
func SubscribeNewBlock() {
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
      record := transformToTxBlock(block)
      log.WithFields(log.Fields{
        "Number":     record.Number,
        "Difficulty": block.Difficulty().String(),
        "Hash":       block.Hash().String(),
      }).Info("A new block generated")
      util.Run(func(db *gorm.DB) error {
        db.NewRecord(record)
        db.Create(record)
        return nil
      })
    }
  }
}

func transformToTxBlock(b *types.Block) *txblock.TxBlock {
  return &txblock.TxBlock{
    Number:      txblock.UInt64(b.Number().Uint64()),
    Hash:        b.Hash(),
    Time:        time.Unix(b.Time().Int64(), 0),
    Txn:         len(b.Transactions()),
    ParentHash:  b.ParentHash(),
    UncleHash:   b.UncleHash(),
    Coinbase:    b.Coinbase(),
    Root:        b.Root(),
    TxHash:      b.TxHash(),
    ReceiptHash: b.ReceiptHash(),
    Bloom:       b.Bloom().Bytes(),
    Difficulty:  txblock.BigInt(*b.Difficulty()),
    GasLimit:    txblock.UInt64(b.GasLimit()),
    GasUsed:     txblock.UInt64(b.GasUsed()),
    Extra:       b.Extra(),
    MixDigest:   b.MixDigest(),
    Nonce:       txblock.UInt64(b.Nonce()),
    // Transactions: func(ts []transactionInfo) []txblock.Transaction {
    //   transactions := make([]txblock.Transaction, len(ts))
    //   for i, t := range b.Transactions {
    //     transactions[i] = txblock.Transaction{
    //       AccountNonce:     util.Parse(t.AccountNonce, util.UInt64Type).(uint64),
    //       Price:            util.Parse(t.Price, util.BytesType).([]byte),
    //       GasLimit:         util.Parse(t.GasLimit, util.UInt64Type).(uint64),
    //       Recipient:        util.Parse(t.Recipient, util.BytesType).([]byte),
    //       Amount:           util.Parse(t.Amount, util.BytesType).([]byte),
    //       Payload:          util.Parse(t.Payload, util.BytesType).([]byte),
    //       V:                util.Parse(t.V, util.BytesType).([]byte),
    //       R:                util.Parse(t.R, util.BytesType).([]byte),
    //       S:                util.Parse(t.S, util.BytesType).([]byte),
    //       Hash:             util.Parse(t.Hash, util.BytesType).([]byte),
    //       BlockNumber:      util.Parse(t.BlockNumber, util.UInt64Type).(uint64),
    //       BlockHash:        util.Parse(t.BlockHash, util.BytesType).([]byte),
    //       From:             util.Parse(t.From, util.BytesType).([]byte),
    //       TransactionIndex: util.Parse(t.TransactionIndex, util.UInt32Type).(uint32),
    //     }
    //   }
    //   return transactions
    // }(b.Transactions),
  }
}
