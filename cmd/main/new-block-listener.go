package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/cypherium/cypherBFT/core/types"
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
	nTxBlock := big.NewInt(0)
	nKeyBlock := big.NewInt(0)
	ticker := time.NewTicker(100 * time.Millisecond)
	newestTicker := time.NewTicker(150 * time.Millisecond)
	blocks := make([]*types.Block, 0, 1000)
	newestBlock := make([]*types.Block, 0, 1)
	latestKeyBlockNumber, err := listerner.BlockFetcher.GetLatestKeyBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}
	latestBlockNumber, err := listerner.BlockFetcher.GetLatestBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}
	nKeyBlock = big.NewInt(latestKeyBlockNumber)
	nTxBlock = big.NewInt(latestBlockNumber)
	log.Infof("latestKeyBlockNumber %d", latestKeyBlockNumber)
	log.Infof("latestBlockNumber %d", latestBlockNumber)
	//log.Infof("localHighestTxBlock time %d", localHighestTxBlock.Time.String())
	log.Infof("nKeyBlock %d", nTxBlock.Uint64())
	log.Infof("nTxBlock %d", nTxBlock.Uint64())
	currentKeyBlock, err := listerner.BlockFetcher.KeyBlockByNumber(big.NewInt(latestKeyBlockNumber))
	if err != nil {
		log.Error("err", fmt.Sprintf("%v", err))
		return
	}
	if currentKeyBlock.Number() == nil {
		currentKeyBlock.SetNumber(big.NewInt(latestKeyBlockNumber))
	}
	for {
		select {
		case newHead := <-newHeader:
			log.Infof("Got new block head time = %v, number = %d, Signature = %x \n\r", time.Unix(0, newHead.Time.Int64()), newHead.Number.Int64(), newHead.Signature)
			block, _, err := listerner.BlockFetcher.BlockByNumber(newHead.Number, true)
			if err != nil {
				log.Error(err)
				return
			}
			//log.Infof("block time", block.Time().String())
			blocks = append(blocks, block)
			latestKeyBlock, err := listerner.BlockFetcher.KeyBlockByNumber(big.NewInt(latestKeyBlockNumber))
			if err != nil {
				log.Error(err)
				return
			}
			if err := listerner.Repo.SaveBlock(block); err == nil {
				newestBlock = append(newestBlock, block)
				listerner.Broadcastable.Broadcast(transformTxBlocksToFrontendMessage(newestBlock, metrics{currentKeyBlock: latestKeyBlock}))
				newestBlock = nil
			}
		case <-ticker.C:
			if nTxBlock.Int64() > 0 {
				block, _, _ := listerner.BlockFetcher.BlockByNumber(nTxBlock, true)
				blocks = append(blocks, block)
				if err := listerner.Repo.SaveBlock(block); err == nil {
					if blocks != nil && len(blocks) > 0 {
						//log.Infof("Broadcst %d blocks", len(blocks))
						listerner.Broadcastable.Broadcast(transformTxBlocksToFrontendMessage(blocks, metrics{currentKeyBlock: currentKeyBlock}))
						blocks = nil
					}
					nTxBlock = nTxBlock.Sub(nTxBlock, big.NewInt(1))
				}
			}
			if nKeyBlock.Int64() > 0 {
				keyBlock, _ := listerner.BlockFetcher.KeyBlockByNumber(nKeyBlock)
				currentKeyBlock = keyBlock
				if err := listerner.Repo.SaveKeyBlock(keyBlock); err == nil {
					listerner.Broadcastable.Broadcast(transformKeyBlockToFrontendMessage(keyBlock.Header()))
					nKeyBlock = nKeyBlock.Sub(nKeyBlock, big.NewInt(1))
				}
			}

		case <-newestTicker.C:
			latestKeyBlockNumber, err := listerner.BlockFetcher.GetLatestKeyBlockNumber()
			if err != nil {
				log.Error(err)
				return
			}
			latestBlockNumber, err := listerner.BlockFetcher.GetLatestBlockNumber()
			if err != nil {
				log.Error(err)
				return
			}
			latestKeyBlock, err := listerner.BlockFetcher.KeyBlockByNumber(big.NewInt(latestKeyBlockNumber))
			if err != nil {
				log.Error(err)
				return
			}
			latestBlock, _, err := listerner.BlockFetcher.BlockByNumber(big.NewInt(latestBlockNumber), true)
			if err != nil {
				log.Error(err)
			}
			if err := listerner.Repo.SaveKeyBlock(latestKeyBlock); err == nil {
				listerner.Broadcastable.Broadcast(transformKeyBlockToFrontendMessage(latestKeyBlock.Header()))
			}
			if err := listerner.Repo.SaveBlock(latestBlock); err == nil {
				newestBlock = append(newestBlock, latestBlock)
				listerner.Broadcastable.Broadcast(transformTxBlocksToFrontendMessage(newestBlock, metrics{currentKeyBlock: latestKeyBlock}))
				newestBlock = nil
			}
		case newKeyHead := <-keyHeadChan:
			log.Infof("Got new kyeBlock timestamp %s\n\r", newKeyHead.Time)
			keyBlock, err := listerner.BlockFetcher.KeyBlockByNumber(newKeyHead.Number)
			if err != nil {
				log.Error(err)
				return
			}
			currentKeyBlock = keyBlock
			log.Infof("Got new key block head: hash = %s, number = %d %v\n\r", newKeyHead.Hash().Hex(), newKeyHead.Number.Int64(), keyBlock.Body().Signatrue)
			if err := listerner.Repo.SaveKeyBlock(keyBlock); err == nil {
				listerner.Broadcastable.Broadcast(transformKeyBlockToFrontendMessage(keyBlock.Header()))
			}

		default:

		}
	}
}
