package repo

import (
	"time"

	"github.com/cypherium/CypherTestNet/go-cypherium/common"
)

// KeyBlock is the Database Table class
type KeyBlock struct {
	Hash   common.Hash `json:"hash"             gencodec:"required"       gorm:"primary_key"`
	Number int64       `json:"number"           gencodec:"required"`
	Time   time.Time   `json:"timestamp"        gencodec:"required"`
	// ParentHash   common.Hash    `json:"parentHash"       gencodec:"required"`
	// Coinbase     common.Address `json:"miner"            gencodec:"required"`
	// Root         common.Hash    `json:"stateRoot"        gencodec:"required"`
	// ReceiptHash  common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	// Bloom        []byte         `json:"logsBloom"        gencodec:"required"`
	// GasLimit     UInt64         `json:"gasLimit"         gencodec:"required"`
	// GasUsed      UInt64         `json:"gasUsed"          gencodec:"required"`
	// Extra        []byte         `json:"extraData"        gencodec:"required"`
	Difficulty BigInt `json:"difficulty"       gencodec:"required"       gorm:"type:blob"`
	// MixDigest    common.Hash    `json:"mixHash"          gencodec:"required"`
	// Nonce        UInt64         `json:"nonce"            gencodec:"required"`
	// UncleHash    common.Hash    `json:"sha3Uncles"       gencodec:"required"`
}
