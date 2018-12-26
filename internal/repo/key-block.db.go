package repo

import (
	"time"

	"github.com/cypherium/CypherTestNet/go-cypherium/common"
)

// KeyBlock is the Database Table class
type KeyBlock struct {
	Hash       Hash        `json:"hash"             gencodec:"required"       gorm:"primary_key"`
	Number     int64       `json:"number"           gencodec:"required"`
	Time       time.Time   `json:"timestamp"        gencodec:"required"`
	ParentHash common.Hash `json:"parentHash"       gencodec:"required"`
	Difficulty BigInt      `json:"difficulty"       gencodec:"required"       gorm:"type:blob"`
	MixDigest  common.Hash `json:"mixHash"          gencodec:"required"`
	Nonce      UInt64      `json:"nonce"            gencodec:"required"`
}
