package blockchain

import (
	"math/big"

	"github.com/cypherium/cypherBFT/core/types"
)

// BlockFetcher is the interface to fetch the Block and Keyblock
type BlockFetcher interface {
	BlockByNumber(number *big.Int, incTx bool) (*types.Block, int, error)
	KeyBlockByNumber(number *big.Int) (*types.KeyBlock, error)
	GetLatestBlockNumber() (int64, error)
	GetLatestKeyBlockNumber() (int64, error)
}

// BlocksFetcher is the interface to fetch the Blocks
type BlocksFetcher interface {
	BlockHeadersByNumbers(numbers []int64) ([]*types.Header, []int, error)
	KeyBlocksByNumbers(numbers []int64) ([]*types.KeyBlock, error)
	GetLatestBlockNumber() (int64, error)
	GetLatestKeyBlockNumber() (int64, error)
}

// Subscription is an interface of subscribe to new block and new key block
type Subscription interface {
	Subscribe(chBlock chan<- *types.Header, chKeyBlock chan<- *types.KeyBlockHeader) (Subscribed, error)
}

// Subscribed is the interface returned by the subscribe func
type Subscribed interface {
	Unsubscribe()
	Err() <-chan error
}
