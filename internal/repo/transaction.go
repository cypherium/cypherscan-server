package repo

import (
	"github.com/ethereum/go-ethereum/common"
)

// Transaction is Transaction struct
type Transaction struct {
	GasPrice         BigInt         `json:"gasPrice" gencodec:"required"       gorm:"type:blob"`
	Gas              UInt64         `json:"gas"      gencodec:"required"`
	To               common.Address `json:"to"`
	Value            BigInt         `json:"value"    gencodec:"required"       gorm:"type:blob"`
	Cost             BigInt         `json:"cost"     gencodec:"required"       gorm:"type:blob"`
	Payload          []byte         `json:"input"    gencodec:"required"`
	Hash             common.Hash    `json:"hash"                               gorm:"primary_key"`
	BlockHash        common.Hash    `json:"blockHash"`
	Block            TxBlock        `json:"-"                                  gorm:"foreignkey:BlockHash"`
	TransactionIndex uint32         `json:"transactionIndex"`

	From common.Address `json:"from"`
	// AccountNonce     UInt64         `json:"nonce"    gencodec:"required"`
	// V                BigInt         `json:"v"        gencodec:"required"       gorm:"type:blob"`
	// R                BigInt         `json:"r"        gencodec:"required"       gorm:"type:blob"`
	// S                BigInt         `json:"s"        gencodec:"required"       gorm:"type:blob"`
}
