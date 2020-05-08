package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/cypherium/cypherBFT/go-cypherium/core/types"
	"github.com/cypherium/cypherscan-server/internal/blockchain"
	"github.com/cypherium/cypherscan-server/internal/publisher"
	"github.com/cypherium/cypherscan-server/internal/repo"
	log "github.com/sirupsen/logrus"
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
	// nTxBlock := big.NewInt(0)
	nKeyBlock := big.NewInt(0)
	ticker := time.NewTicker(2 * time.Second)
	blocks := make([]*types.Block, 0, 1000)
	latestKeyBlocksNumber, _ := listerner.BlockFetcher.GetLatestKeyBlockNumber()
	currentKeyBlock, err := listerner.BlockFetcher.KeyBlockByNumber(big.NewInt(latestKeyBlocksNumber))
	if err != nil {
		log.Error("err", fmt.Sprintf("%v", err))
		return
	}
	if currentKeyBlock.Number() == nil {
		currentKeyBlock.SetNumber(big.NewInt(latestKeyBlocksNumber))
	}
	listerner.Broadcastable.Broadcast(transformTxBlocksToFrontendMessage([]*types.Block{}, metrics{currentKeyBlock: currentKeyBlock}))

	// _k, err := listerner.BlockFetcher.KeyBlockByNumber(big.NewInt(400))
	// if err != nil {
	// 	fmt.Printf("ERrrrrrror, %s", err.Error())
	// } else {
	// 	fmt.Printf("got b: %x %x\n", _k.Body().LeaderPubKey, _k.Body().Signatrue)
	// }

	for {
		select {

		case newHead := <-newHeader:
			fmt.Printf("Got new block head time = %v, number = %d, Signature = %x \n\r", time.Unix(0, newHead.Time.Int64()), newHead.Number.Int64(), newHead.Signature)
			block, _, _ := listerner.BlockFetcher.BlockByNumber(newHead.Number, true)
			blocks = append(blocks, block)
			listerner.Repo.SaveBlock(block)
			listerner.BlockFetcher.SetLatestNumbers(newHead.Number.Int64(), -1)
		case <-ticker.C:
			if blocks != nil && len(blocks) > 0 {
				fmt.Printf("Broadcst %d blocks", len(blocks))
				listerner.Broadcastable.Broadcast(transformTxBlocksToFrontendMessage(blocks, metrics{currentKeyBlock: currentKeyBlock}))
				blocks = nil
			}
		case newKeyHead := <-keyHeadChan:
			keyBlock, _ := listerner.BlockFetcher.KeyBlockByNumber(newKeyHead.Number)
			currentKeyBlock = keyBlock
			fmt.Printf("Got new key block head: hash = %s, number = %d %v\n\r", newKeyHead.Hash().Hex(), newKeyHead.Number.Int64(), keyBlock.Body().Signatrue)
			listerner.Repo.SaveKeyBlock(keyBlock)
			listerner.Broadcastable.Broadcast(transformKeyBlockToFrontendMessage(newKeyHead))
			listerner.BlockFetcher.SetLatestNumbers(-1, newKeyHead.Number.Int64())

		default:
			// latestBlocksNumber, _ := listerner.BlockFetcher.GetLatestBlockNumber()
			latestKeyBlocksNumber, _ := listerner.BlockFetcher.GetLatestKeyBlockNumber()
			// if nTxBlock.Int64() <= latestBlocksNumber {
			// 	block, _, _ := listerner.BlockFetcher.BlockByNumber(nTxBlock, true)
			// 	blocks = append(blocks, block)
			// 	listerner.Repo.SaveBlock(block)
			// 	listerner.BlockFetcher.SetLatestNumbers(nTxBlock.Int64(), -1)
			// 	nTxBlock = nTxBlock.Add(nTxBlock, big.NewInt(1))
			// }
			if nKeyBlock.Int64() <= latestKeyBlocksNumber {
				keyBlock, _ := listerner.BlockFetcher.KeyBlockByNumber(nKeyBlock)
				currentKeyBlock = keyBlock
				listerner.Repo.SaveKeyBlock(keyBlock)
				listerner.Broadcastable.Broadcast(transformKeyBlockToFrontendMessage(keyBlock.Header()))
				// listerner.BlockFetcher.SetLatestNumbers(-1, nKeyBlock.Int64())
				nKeyBlock = nKeyBlock.Add(nKeyBlock, big.NewInt(1))
			}
			time.Sleep(500 * time.Millisecond)

		}
	}
}
