package txblock

import (
	"math/big"
)

const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 20
	// BloomByteLength represents the number of bytes used in a header log bloom.
	BloomByteLength = 256
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// Address represents the 20 byte address of an Ethereum account.
type Address [AddressLength]byte

// Bloom represents a 2048 bit bloom filter.
type Bloom [BloomByteLength]byte

// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type BlockNonce [8]byte

// TxBlock is the Database Table class
type TxBlock struct {
	ParentHash  Hash       `json:"parentHash"       gencodec:"required"`
	UncleHash   Hash       `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    Address    `json:"miner"            gencodec:"required"`
	Root        Hash       `json:"stateRoot"        gencodec:"required"`
	TxHash      Hash       `json:"transactionsRoot" gencodec:"required"     gorm:"primary_key"`
	ReceiptHash Hash       `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom      `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int   `json:"difficulty"       gencodec:"required"`
	Number      *big.Int   `json:"number"           gencodec:"required"`
	GasLimit    uint64     `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64     `json:"gasUsed"          gencodec:"required"`
	Time        *big.Int   `json:"timestamp"        gencodec:"required"`
	Extra       []byte     `json:"extraData"        gencodec:"required"`
	MixDigest   Hash       `json:"mixHash"          gencodec:"required"`
	Nonce       BlockNonce `json:"nonce"            gencodec:"required"`
}
