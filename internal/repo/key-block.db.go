package repo

import (
	"time"
)

// KeyBlock is the Database Table class
type KeyBlock struct {
	Hash       Hash      `json:"hash"             gencodec:"required"       gorm:"primary_key"`
	Number     int64     `json:"number"           gencodec:"required"`
	Time       time.Time `json:"timestamp"        gencodec:"required"`
	ParentHash Hash      `json:"parentHash"       gencodec:"required"`
	Difficulty BigInt    `json:"difficulty"       gencodec:"required"`
	MixDigest  Hash      `json:"mixHash"          gencodec:"required"`
	Nonce      UInt64    `json:"nonce"            gencodec:"required"`
}
