package blockchain

import ethereum "github.com/cypherium/cypherBFT/go-cypherium/event"

// CypherSubscribed is the struct subscribed to blockchain
type CypherSubscribed struct {
	BlockSubscription    ethereum.Subscription
	KeyBlockSubscription ethereum.Subscription
}

// Unsubscribe is to unsubscribe
func (subscribed *CypherSubscribed) Unsubscribe() {
	subscribed.BlockSubscription.Unsubscribe()
	subscribed.KeyBlockSubscription.Unsubscribe()
}

// Err is to merge to Blockchains' err
func (subscribed *CypherSubscribed) Err() <-chan error {
	return nil
}
