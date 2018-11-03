package txblock

import (
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
  "math/big"
)

// Transaction is Transaction struct
type Transaction struct {
  AccountNonce     uint64       `json:"nonce"    gencodec:"required"`
  Price            big.Int      `json:"gasPrice" gencodec:"required"`
  GasLimit         uint64       `json:"gas"      gencodec:"required"`
  Recipient        util.Address `json:"to"       rlp:"nil"` // nil means contract creation
  Amount           big.Int      `json:"value"    gencodec:"required"`
  Payload          []byte       `json:"input"    gencodec:"required"`
  V                big.Int      `json:"v"        gencodec:"required"`
  R                big.Int      `json:"r"        gencodec:"required"`
  S                big.Int      `json:"s"        gencodec:"required"`
  Hash             util.Hash    `json:"hash"     rlp:"-"`
  BlockNumber      TxBlock      `json:"blockNumber"                         gorm:"foreignkey:Number"`
  BlockHash        util.Hash    `json:"blockHash"`
  From             util.Address `json:"from"`
  TransactionIndex uint         `json:"transactionIndex"`
}
