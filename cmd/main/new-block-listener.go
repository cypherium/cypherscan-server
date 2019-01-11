package main

import (
	"fmt"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	"gitlab.com/ron-liu/cypherscan-server/internal/blockchain"
	"gitlab.com/ron-liu/cypherscan-server/internal/publisher"

	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
)

// NewBlockListener is to listen the
type NewBlockListener struct {
	Repo          repo.BlockSaver
	BlockFetcher  blockchain.BlockFetcher
	Broadcastable publisher.Broadcastable
}

// Listen is to listen
func (listerner *NewBlockListener) Listen(newHeader chan *types.Header, keyHeadChan chan *types.KeyBlockHeader) {
	for {
		select {
		case newHead := <-newHeader:
			fmt.Printf("Got new block head hash = %s, number = %d \n\r", newHead.Hash().Hex(), newHead.Number.Int64())
			block, _, _ := listerner.BlockFetcher.BlockByNumber(newHead.Number, true)
			listerner.Repo.SaveBlock(block)
			listerner.Broadcastable.Broadcast(transformTxBlockToFrontendMessage(block))
			listerner.BlockFetcher.SetLatestNumbers(newHead.Number.Int64(), -1)

		case newKeyHead := <-keyHeadChan:
			fmt.Printf("Got new key block head: hash = %s, number = %d\n\r", newKeyHead.Hash().Hex(), newKeyHead.Number.Int64())
			listerner.Repo.SaveKeyBlock(newKeyHead)
			listerner.Broadcastable.Broadcast(transformKeyBlockToFrontendMessage(newKeyHead))
			listerner.BlockFetcher.SetLatestNumbers(-1, newKeyHead.Number.Int64())
		}
	}
}