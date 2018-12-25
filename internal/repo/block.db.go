package repo

import (
	"time"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
)

// TxBlock is the Database Table class
type TxBlock struct {
	Hash         Hash          `json:"hash"             gencodec:"required"       gorm:"primary_key"`
	Number       int64         `json:"number"           gencodec:"required"`
	Time         time.Time     `json:"timestamp"        gencodec:"required"`
	Txn          int           `json:"txn"              gencodec:"required"`
	ParentHash   Hash          `json:"parentHash"       gencodec:"required"`
	Coinbase     Address       `json:"miner"            gencodec:"required"`
	Root         Hash          `json:"stateRoot"        gencodec:"required"`
	TxHash       Hash          `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash  Hash          `json:"receiptsRoot"     gencodec:"required"`
	Bloom        []byte        `json:"logsBloom"        gencodec:"required"`
	GasLimit     UInt64        `json:"gasLimit"         gencodec:"required"`
	GasUsed      UInt64        `json:"gasUsed"          gencodec:"required"`
	Extra        []byte        `json:"extraData"        gencodec:"required"`
	Transactions []Transaction `json:"transactions"     gencodec:"required"     gorm:"foreignkey:BlockHash"`
	Difficulty   BigInt        `json:"difficulty"       gencodec:"required"       gorm:"type:blob"`
	// MixDigest    common.Hash    `json:"mixHash"          gencodec:"required"`
	// Nonce        UInt64         `json:"nonce"            gencodec:"required"`
	// UncleHash    common.Hash    `json:"sha3Uncles"       gencodec:"required"`
}

func transformBlockToDbRecord(b *types.Block) *TxBlock {
	return &TxBlock{
		Number:      b.Number().Int64(),
		Hash:        Hash(b.Hash()),
		Time:        time.Unix(b.Time().Int64(), 0),
		Txn:         len(b.Transactions()),
		ParentHash:  Hash(b.ParentHash()),
		Coinbase:    Address(b.Coinbase()),
		Root:        Hash(b.Root()),
		TxHash:      Hash(b.TxHash()),
		ReceiptHash: Hash(b.ReceiptHash()),
		Bloom:       b.Bloom().Bytes(),
		// Difficulty:  BigInt(*b.Difficulty()),
		GasLimit: UInt64(b.GasLimit()),
		GasUsed:  UInt64(b.GasUsed()),
		Extra:    b.Extra(),
		// UncleHash:   b.UncleHash(),
		// MixDigest:   b.MixDigest(),
		// Nonce:       UInt64(b.Nonce()),
		Transactions: func(ts []*types.Transaction) []Transaction {
			transactions := make([]Transaction, len(ts))
			for i, t := range ts {
				transactions[i] = Transaction{
					Hash:     Hash(t.Hash()),
					Gas:      UInt64(t.Gas()),
					GasPrice: BigInt(*t.GasPrice()),
					To: func() Address {
						if t.To() != nil {
							return Address(*t.To())
						}
						return Address{}
					}(),
					Value:            BigInt(*t.Value()),
					Cost:             BigInt(*t.Cost()),
					BlockHash:        Hash(b.Hash()),
					TransactionIndex: uint32(i),
					Payload:          t.Data(),
					// Recipient:        util.Parse(t.Recipient, util.BytesType).([]byte),
					// AccountNonce:     UInt64(t.Nonce()),
					// V: func() BigInt {
					//   v, _, _ := t.RawSignatureValues()
					//   return BigInt(*v)
					// }(),
					// R: func() BigInt {
					//   _, r, _ := t.RawSignatureValues()
					//   return BigInt(*r)
					// }(),
					// S: func() BigInt {
					//   _, _, s := t.RawSignatureValues()
					//   return BigInt(*s)
					// }(),
				}
			}
			return transactions
		}(b.Transactions()),
	}
}
