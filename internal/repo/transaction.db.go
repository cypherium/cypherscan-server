package repo

// Transaction is Transaction struct
type Transaction struct {
	Hash             Hash    `json:"hash"                               gorm:"primary_key"`
	GasPrice         BigInt  `json:"gasPrice" gencodec:"required"       gorm:"type:blob"`
	Gas              UInt64  `json:"gas"      gencodec:"required"`
	From             Address `json:"from"`
	To               Address `json:"to"`
	Value            BigInt  `json:"value"    gencodec:"required"       gorm:"type:blob"`
	Cost             BigInt  `json:"cost"     gencodec:"required"       gorm:"type:blob"`
	Payload          []byte  `json:"input"    gencodec:"required"`
	TransactionIndex uint32  `json:"transactionIndex"`
	BlockHash        Hash    `json:"blockHash"`
	BlockNumber      int64   `json:"blockNumber"`
	Block            TxBlock `json:"-"                                  gorm:"foreignkey:BlockNumber"`
}
