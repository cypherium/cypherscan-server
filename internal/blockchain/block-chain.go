package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	"github.com/cypherium/CypherTestNet/go-cypherium/ethclient"
	"gitlab.com/ron-liu/cypherscan-server/internal/util"
)

// BlockChain is the struct of the Client
type BlockChain struct {
	client               *ethclient.Client
	context              context.Context
	latestBlockNumber    int64
	latestKeyBlockNumber int64
}

// BlockHeadersByNumbers is to get BlockHeaders by numbers
func (blockChain *BlockChain) BlockHeadersByNumbers(numbers []int64) ([]*types.Header, error) {
	// return blockChain.client.BlockHeadersByNumbers(blockChain.context, numbers)
	headers, error := blockChain.client.BlockHeadersByNumbers(blockChain.context, numbers)
	if error == nil {
		for _, h := range headers {
			setToCurrentTime(h, nil)
		}
	}
	return headers, error
}

// BlockByNumber is to get Block by number, if number is nil , and return number of transactions without retreive whole transation slice
func (blockChain *BlockChain) BlockByNumber(number *big.Int, incTx bool) (*types.Block, int, error) {
	// return blockChain.client.BlockByNumber(blockChain.context, number, incTx)
	block, txn, err := blockChain.client.BlockByNumber(blockChain.context, number, incTx)
	if err == nil {
		setToCurrentTime(nil, block)
	}
	return block, txn, err
}

// KeyBlockByNumber is to get Key Block by number
func (blockChain *BlockChain) KeyBlockByNumber(number *big.Int) (*types.KeyBlock, error) {
	return blockChain.client.KeyBlockByNumber(blockChain.context, number)
}

// KeyBlocksByNumbers is to get BlockHeaders by numbers
func (blockChain *BlockChain) KeyBlocksByNumbers(numbers []int64) ([]*types.KeyBlock, error) {
	return blockChain.client.KeyBlocksByNumbers(blockChain.context, numbers)
}

// GetLatestBlockNumber is to get the latest Block number
func (blockChain *BlockChain) GetLatestBlockNumber() (int64, error) {
	if blockChain.latestBlockNumber <= 0 {
		b, _, err := blockChain.client.BlockByNumber(blockChain.context, nil, false)
		if err != nil {
			return 0, err
		}
		blockChain.latestBlockNumber = b.Number().Int64()
	}
	return blockChain.latestBlockNumber, nil
}

// GetLatestKeyBlockNumber is to get the latest KeyBlock Number
func (blockChain *BlockChain) GetLatestKeyBlockNumber() (int64, error) {
	if blockChain.latestKeyBlockNumber <= 0 {
		b, err := blockChain.client.KeyBlockByNumber(blockChain.context, nil)
		if err != nil {
			fmt.Printf("xxxxxx: %s\n", err.Error())
			return 0, err
		}
		blockChain.latestKeyBlockNumber = b.Number().Int64()
	}
	return blockChain.latestKeyBlockNumber, nil
}

// SetLatestNumbers is to set the latest block/key block number
func (blockChain *BlockChain) SetLatestNumbers(blockNumber int64, keyBlockNumber int64) {
	if blockNumber > 0 {
		blockChain.latestBlockNumber = blockNumber
	}
	if keyBlockNumber > 0 {
		blockChain.latestKeyBlockNumber = keyBlockNumber
	}
}

// Subscribe is to subscirbe new block and new key block
func (blockChain *BlockChain) Subscribe(chBlock chan<- *types.Header, chKeyBlock chan<- *types.KeyBlockHeader) (Subscribed, error) {
	blockSub, err := blockChain.client.SubscribeNewHead(blockChain.context, chBlock)
	if err != nil {
		log.Error(err.Error())
		return nil, &util.MyError{Message: "Cannot subscirbe to Block"}
	}
	keyBlockSub, err := blockChain.client.SubscribeNewKeyHead(blockChain.context, chKeyBlock)
	if err != nil {
		log.Error(err.Error())
		blockSub.Unsubscribe()
		return nil, &util.MyError{Message: "Cannot subscirbe to Key Block"}
	}

	return &CypherSubscribed{blockSub, keyBlockSub}, nil
}

func setToCurrentTime(header *types.Header, block *types.Block) {
	if header != nil {
		header.Time = big.NewInt(time.Now().UnixNano())
	}
	if block != nil {
		block.SetToCurrentTime()
	}
}
