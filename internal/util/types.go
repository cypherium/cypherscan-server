package util

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
