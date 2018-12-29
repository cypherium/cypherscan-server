package blockchain

import (
	"github.com/cypherium/CypherTestNet/go-cypherium/common"
	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
)

// BlockFetcher is the interface to fetch the Block and Keyblock
type BlockFetcher interface {
	BlockByHash(hash common.Hash, incTx bool) (*types.Block, int, error)
	KeyBlockByHash(hash common.Hash) (*types.KeyBlock, error)
}

//BlocksFetcher is the interface to get multiple blocks
type BlocksFetcher interface {
	BlockHeadersByNumbers(numbers []int64) ([]*types.Header, error)
	KeyBlocksByNumbers(numbers []int64) ([]*types.KeyBlock, error)
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
