package repo

import (
	"time"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	"github.com/cypherium/CypherTestNet/go-cypherium/crypto"
)

// TxBlock is the Database Table class
type TxBlock struct {
	ID           int64         `json:"-" gorm:"primary_key"`
	Number       int64         `json:"number"`
	Hash         Hash          `json:"hash"`
	Time         time.Time     `json:"timestamp"`
	Txn          int           `json:"txn"`
	ParentHash   Hash          `json:"parentHash"`
	Root         Hash          `json:"stateRoot"`
	TxHash       Hash          `json:"transactionsRoot"`
	ReceiptHash  Hash          `json:"receiptsRoot"`
	Bloom        []byte        `json:"logsBloom"`
	GasLimit     UInt64        `json:"gasLimit"`
	GasUsed      UInt64        `json:"gasUsed"`
	Transactions []Transaction `json:"transactions" gorm:"foreignkey:BlockNumber;association_foreignkey:Number"`
	KeySignature Bytes         `json:"keySignature"`
}

func transferKeyBlockHeaderToDbRecord(b *types.KeyBlockHeader) *KeyBlock {
	return &KeyBlock{
		Hash:       Hash(b.Hash()),
		Number:     b.Number.Int64(),
		Difficulty: UInt64(b.Difficulty.Uint64()),
		Time:       time.Unix(b.Time.Int64(), 0),
	}
}

func transformBlockToDbRecord(b *types.Block) *TxBlock {
	return &TxBlock{
		Number:       b.Number().Int64(),
		Hash:         Hash(b.Hash()),
		Time:         time.Unix(0, b.Time().Int64()),
		Txn:          len(b.Transactions()),
		ParentHash:   Hash(b.ParentHash()),
		Root:         Hash(b.Root()),
		TxHash:       Hash(b.TxHash()),
		ReceiptHash:  Hash(b.ReceiptHash()),
		Bloom:        b.Bloom().Bytes(),
		GasLimit:     UInt64(b.GasLimit()),
		GasUsed:      UInt64(b.GasUsed()),
		KeySignature: Bytes(b.Header().KeySignature),
		Transactions: func(ts []*types.Transaction) []Transaction {
			transactions := make([]Transaction, len(ts))
			for i, t := range ts {
				transactions[i] = Transaction{
					Hash:     Hash(t.Hash()),
					Gas:      UInt64(t.Gas()),
					GasPrice: UInt64(t.GasPrice().Uint64()),
					To: func() Address {
						if t.To() != nil {
							return Address(*t.To())
						}
						return Address{}
					}(),
					From:             Address(crypto.PubKeyToAddressCypherium(t.SenderKey())),
					Value:            UInt64(t.Value().Uint64()),
					Cost:             UInt64(t.Cost().Uint64()),
					BlockHash:        Hash(b.Hash()),
					BlockNumber:      b.Number().Int64(),
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
