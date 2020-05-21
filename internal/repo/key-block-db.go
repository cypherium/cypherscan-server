package repo

import (
	"time"

	"github.com/cypherium/cypherBFT/go-cypherium/core/types"
)

// KeyBlock is the Database Table class
type KeyBlock struct {
	Hash         Hash      `json:"hash" gorm:"primary_key"`
	Number       int64     `json:"number" gorm:"indßex:key_blocks_number`
	Time         time.Time `json:"timestamp"`
	ParentHash   Hash      `json:"parentHash"`
	Difficulty   UInt64    `json:"difficulty"`
	MixDigest    Hash      `json:"mixHash"`
	Nonce        UInt64    `json:"nonce"`
	Signature    Bytes     `json:"signature"`
	LeaderPubKey Bytes     `json:"leaderPubKey"`
}

func transferKeyBlockHeaderToDbRecord(b *types.KeyBlock) *KeyBlock {
	return &KeyBlock{
		Hash:         Hash(b.Hash()),
		Number:       b.Number().Int64(),
		Difficulty:   UInt64(b.Difficulty().Uint64()),
		Time:         time.Unix(b.Time().Int64(), 0),
		Signature:    Bytes(b.Body().Signatrue),
		LeaderPubKey: Bytes(b.Body().LeaderPubKey),
		Nonce:        UInt64(b.Nonce()),
		MixDigest:    Hash(b.MixDigest()),
	}
}
