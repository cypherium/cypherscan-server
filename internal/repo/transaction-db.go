package repo

// Transaction is Transaction struct
type Transaction struct {
	ID               int64   `json:"-" gorm:"primary_key"`
	Hash             Hash    `json:"hash" gorm:"index:transactions_hash`
	GasPrice         UInt64  `json:"gasPrice"`
	Gas              UInt64  `json:"gas"`
	From             Address `json:"from"`
	To               Address `json:"to"`
	Value            UInt64  `json:"value"`
	Cost             UInt64  `json:"cost"`
	Payload          []byte  `json:"input"`
	TransactionIndex uint32  `json:"transactionIndex"`
	BlockHash        Hash    `json:"blockHash"`
	BlockNumber      int64   `json:"blockNumber"`
	Block            TxBlock `json:"block" gorm:"foreignkey:BlockNumber;association_foreignkey:Number"`
	Signature        Bytes   `json:"signature" `
}
