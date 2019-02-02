package main

import (
	"fmt"
	"time"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	"gitlab.com/ron-liu/cypherscan-server/internal/blockchain"
	"gitlab.com/ron-liu/cypherscan-server/internal/publisher"

	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
)

const (
	BroadcastInterval = 5
)

// NewBlockListener is to listen the
type NewBlockListener struct {
	Repo          repo.BlockSaver
	BlockFetcher  blockchain.BlockFetcher
	Broadcastable publisher.Broadcastable
}

// Listen is to listen
func (listerner *NewBlockListener) Listen(newHeader chan *types.Header, keyHeadChan chan *types.KeyBlockHeader) {
	ticker := time.NewTicker(2 * time.Second)
	blocks := make([]*types.Block, 0, 1000)
	for {
		select {

		case newHead := <-newHeader:
			fmt.Printf("Got new block head time = %v, number = %d, KeySignature = %x \n\r", time.Unix(0, newHead.Time.Int64()), newHead.Number.Int64(), newHead.KeySignature)
			block, _, _ := listerner.BlockFetcher.BlockByNumber(newHead.Number, true)
			blocks = append(blocks, block)
			listerner.Repo.SaveBlock(block)
			listerner.BlockFetcher.SetLatestNumbers(newHead.Number.Int64(), -1)
		case <-ticker.C:
			if blocks != nil || len(blocks) > 0 {
				fmt.Printf("Broadcst %d blocks", len(blocks))
				listerner.Broadcastable.Broadcast(transformTxBlocksToFrontendMessage(blocks))
				blocks = nil
			}
		case newKeyHead := <-keyHeadChan:
			fmt.Printf("Got new key block head: hash = %s, number = %d\n\r", newKeyHead.Hash().Hex(), newKeyHead.Number.Int64())
			listerner.Repo.SaveKeyBlock(newKeyHead)
			listerner.Broadcastable.Broadcast(transformKeyBlockToFrontendMessage(newKeyHead))
			listerner.BlockFetcher.SetLatestNumbers(-1, newKeyHead.Number.Int64())
		}
	}
}
