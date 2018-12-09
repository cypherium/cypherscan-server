package txblock

import (
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/jinzhu/gorm"
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
  "time"
)

// TxBlock is the Database Table class
type TxBlock struct {
  Hash         common.Hash    `json:"hash"             gencodec:"required"       gorm:"primary_key"`
  Number       UInt64         `json:"number"           gencodec:"required"`
  Time         time.Time      `json:"timestamp"        gencodec:"required"`
  Txn          int            `json:"txn"              gencodec:"required"`
  ParentHash   common.Hash    `json:"parentHash"       gencodec:"required"`
  Coinbase     common.Address `json:"miner"            gencodec:"required"`
  Root         common.Hash    `json:"stateRoot"        gencodec:"required"`
  TxHash       common.Hash    `json:"transactionsRoot" gencodec:"required"`
  ReceiptHash  common.Hash    `json:"receiptsRoot"     gencodec:"required"`
  Bloom        []byte         `json:"logsBloom"        gencodec:"required"`
  GasLimit     UInt64         `json:"gasLimit"         gencodec:"required"`
  GasUsed      UInt64         `json:"gasUsed"          gencodec:"required"`
  Extra        []byte         `json:"extraData"        gencodec:"required"`
  Transactions []Transaction  `json:"transactions"     gencodec:"required"     gorm:"foreignkey:BlockHash"`
  Difficulty   BigInt         `json:"difficulty"       gencodec:"required"       gorm:"type:blob"`
  // MixDigest    common.Hash    `json:"mixHash"          gencodec:"required"`
  // Nonce        UInt64         `json:"nonce"            gencodec:"required"`
  // UncleHash    common.Hash    `json:"sha3Uncles"       gencodec:"required"`
}

// SaveBlock to the database
func SaveBlock(block *types.Block) error {
  record := transformBlockToDbRecord(block)
  util.RunDb(func(db *gorm.DB) error {
    db.NewRecord(record)
    db.Create(record)
    return nil
  })
  return nil
}

func transformBlockToDbRecord(b *types.Block) *TxBlock {
  return &TxBlock{
    Number:      UInt64(b.Number().Uint64()),
    Hash:        b.Hash(),
    Time:        time.Unix(b.Time().Int64(), 0),
    Txn:         len(b.Transactions()),
    ParentHash:  b.ParentHash(),
    Coinbase:    b.Coinbase(),
    Root:        b.Root(),
    TxHash:      b.TxHash(),
    ReceiptHash: b.ReceiptHash(),
    Bloom:       b.Bloom().Bytes(),
    Difficulty:  BigInt(*b.Difficulty()),
    GasLimit:    UInt64(b.GasLimit()),
    GasUsed:     UInt64(b.GasUsed()),
    Extra:       b.Extra(),
    // UncleHash:   b.UncleHash(),
    // MixDigest:   b.MixDigest(),
    // Nonce:       UInt64(b.Nonce()),
    Transactions: func(ts []*types.Transaction) []Transaction {
      transactions := make([]Transaction, len(ts))
      for i, t := range ts {
        transactions[i] = Transaction{
          Hash:     t.Hash(),
          Gas:      UInt64(t.Gas()),
          GasPrice: BigInt(*t.GasPrice()),
          To: func() common.Address {
            if t.To() != nil {
              return *(t.To())
            }
            return common.Address{}
          }(),
          Value:            BigInt(*t.Value()),
          Cost:             BigInt(*t.Cost()),
          BlockHash:        b.Hash(),
          TransactionIndex: uint32(i),
          Payload:          t.Data(),
          // Recipient:        util.Parse(t.Recipient, util.BytesType).([]byte),
          // AccountNonce:     UInt64(t.Nonce()),
          // V: func() BigInt {
          //   v, _, _ := t.RawSignatureValues()
          //   return BigInt(*v)
          // }(),
          // R: func() BigInt {
          //   _, r, _ := t.RawSignatureValues()
          //   return BigInt(*r)
          // }(),
          // S: func() BigInt {
          //   _, _, s := t.RawSignatureValues()
          //   return BigInt(*s)
          // }(),
        }
      }
      return transactions
    }(b.Transactions()),
  }
}
