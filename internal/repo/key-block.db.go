package repo

import (
	"time"
)

// KeyBlock is the Database Table class
type KeyBlock struct {
	Hash       Hash      `json:"hash"	        gorm:"primary_key"`
	Number     int64     `json:"number"`
	Time       time.Time `json:"timestamp"`
	ParentHash Hash      `json:"parentHash"`
	Difficulty BigInt    `json:"difficulty"`
	MixDigest  Hash      `json:"mixHash"`
	Nonce      UInt64    `json:"nonce"`
}
