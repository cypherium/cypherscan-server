package txblock

import (
  "github.com/ethereum/go-ethereum/common"
  "time"
)

// TxBlock is the Database Table class
type TxBlock struct {
  Hash        common.Hash    `json:"hash"             gencodec:"required"       gorm:"primary_key"`
  Number      UInt64         `json:"number"           gencodec:"required"       gorm:"type:bigint"`
  Time        time.Time      `json:"timestamp"        gencodec:"required"`
  Txn         int            `json:"txn"              gencodec:"required"       gorm:"type:integer"`
  ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
  UncleHash   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
  Coinbase    common.Address `json:"miner"            gencodec:"required"`
  Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
  TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
  ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
  Bloom       []byte         `json:"logsBloom"        gencodec:"required"`
  Difficulty  BigInt         `json:"difficulty"       gencodec:"required"       gorm:"type:blob"`
  GasLimit    UInt64         `json:"gasLimit"         gencodec:"required"`
  GasUsed     UInt64         `json:"gasUsed"          gencodec:"required"        `
  Extra       []byte         `json:"extraData"        gencodec:"required"`
  MixDigest   common.Hash    `json:"mixHash"          gencodec:"required"`
  Nonce       UInt64         `json:"nonce"            gencodec:"required"`
  // Transactions []Transaction `json:"transactions"     gencodec:"required"     gorm:"foreignkey:BlockNumber"`
}
