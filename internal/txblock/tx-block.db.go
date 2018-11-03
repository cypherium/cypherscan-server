package txblock

import (
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
  "math/big"
  "time"
)

// TxBlock is the Database Table class
type TxBlock struct {
  ParentHash   util.Hash       `json:"parentHash"       gencodec:"required"`
  UncleHash    util.Hash       `json:"sha3Uncles"       gencodec:"required"`
  Coinbase     util.Address    `json:"miner"            gencodec:"required"`
  Root         util.Hash       `json:"stateRoot"        gencodec:"required"`
  TxHash       util.Hash       `json:"transactionsRoot" gencodec:"required"`
  ReceiptHash  util.Hash       `json:"receiptsRoot"     gencodec:"required"`
  Bloom        util.Bloom      `json:"logsBloom"        gencodec:"required"`
  Difficulty   big.Int         `json:"difficulty"       gencodec:"required"`
  Number       uint64          `json:"number"           gencodec:"required"        gorm:"primary_key;type:bigint"`
  GasLimit     uint64          `json:"gasLimit"         gencodec:"required"        gorm:"type:bigint"`
  GasUsed      uint64          `json:"gasUsed"          gencodec:"required"        gorm:"type:bigint"`
  Time         time.Time       `json:"timestamp"        gencodec:"required"`
  Extra        []byte          `json:"extraData"        gencodec:"required"`
  MixDigest    util.Hash       `json:"mixHash"          gencodec:"required"`
  Nonce        util.BlockNonce `json:"nonce"            gencodec:"required"`
  Transactions []Transaction   `json:"transactions"     gencodec:"required"     gorm:"foreignkey:Number"`
}
