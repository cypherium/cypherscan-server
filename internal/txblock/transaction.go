package txblock

// Transaction is Transaction struct
type Transaction struct {
  AccountNonce     uint64 `json:"nonce"    gencodec:"required"`
  Price            []byte `json:"gasPrice" gencodec:"required"`
  GasLimit         uint64 `json:"gas"      gencodec:"required"`
  Recipient        []byte `json:"to"       rlp:"nil"` // nil means contract creation
  Amount           []byte `json:"value"    gencodec:"required"`
  Payload          []byte `json:"input"    gencodec:"required"`
  V                []byte `json:"v"        gencodec:"required"`
  R                []byte `json:"r"        gencodec:"required"`
  S                []byte `json:"s"        gencodec:"required"`
  Hash             []byte `json:"hash"     rlp:"-"                     gorm:"primary_key:Hash"`
  BlockNumber      uint64 `json:"blockNumber"`
  BlockHash        []byte `json:"blockHash"`
  From             []byte `json:"from"`
  TransactionIndex uint32 `json:"transactionIndex"`
}
