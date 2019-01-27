package repo

// Transaction is Transaction struct
type Transaction struct {
	ID               int64   `json:"-" gorm:"primary_key"`
	Hash             Hash    `json:"hash"`
	GasPrice         BigInt  `json:"gasPrice" gencodec:"required"`
	Gas              UInt64  `json:"gas"      gencodec:"required"`
	From             Address `json:"from"`
	To               Address `json:"to"`
	Value            BigInt  `json:"value"    gencodec:"required"`
	Cost             BigInt  `json:"cost"     gencodec:"required"`
	Payload          []byte  `json:"input"    gencodec:"required"`
	TransactionIndex uint32  `json:"transactionIndex"`
	BlockHash        Hash    `json:"blockHash"`
	BlockNumber      int64   `json:"blockNumber"`
	Block            TxBlock `json:"-"                                  gorm:"foreignkey:BlockNumber"`
}
