package blockchain

import "github.com/cypherium/CypherTestNet/go-cypherium/core/types"

// Subscription is an interface of subscribe to new block and new key block
type Subscription interface {
	Subscribe(chBlock chan<- *types.Header, chKeyBlock chan<- *types.KeyBlockHeader)
}

// Subscribed is the interface
type Subscribed interface {
	Unsubscribe()
	Err() <-chan error
}
