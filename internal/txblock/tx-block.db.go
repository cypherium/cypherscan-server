package txblock

import (
  // "gitlab.com/ron-liu/cypherscan-server/internal/util"
  // "math/big"
  "time"
)

// TxBlock is the Database Table class
type TxBlock struct {
  ParentHash   []byte        `json:"parentHash"       gencodec:"required"`
  UncleHash    []byte        `json:"sha3Uncles"       gencodec:"required"`
  Coinbase     []byte        `json:"miner"            gencodec:"required"`
  Root         []byte        `json:"stateRoot"        gencodec:"required"`
  TxHash       []byte        `json:"transactionsRoot" gencodec:"required"`
  ReceiptHash  []byte        `json:"receiptsRoot"     gencodec:"required"`
  Bloom        []byte        `json:"logsBloom"        gencodec:"required"`
  Difficulty   []byte        `json:"difficulty"       gencodec:"required"`
  Number       uint64        `json:"number"           gencodec:"required"        gorm:"primary_key;type:bigint"`
  GasLimit     []byte        `json:"gasLimit"         gencodec:"required"        `
  GasUsed      []byte        `json:"gasUsed"          gencodec:"required"        `
  Time         time.Time     `json:"timestamp"        gencodec:"required"`
  Extra        []byte        `json:"extraData"        gencodec:"required"`
  MixDigest    []byte        `json:"mixHash"          gencodec:"required"`
  Nonce        []byte        `json:"nonce"            gencodec:"required"`
  Transactions []Transaction `json:"transactions"     gencodec:"required"     gorm:"foreignkey:BlockNumber"`
}
