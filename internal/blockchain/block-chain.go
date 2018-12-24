package blockchain

import (
	"context"

	"github.com/cypherium/CypherTestNet/go-cypherium/common"
	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	"github.com/cypherium/CypherTestNet/go-cypherium/ethclient"
	"gitlab.com/ron-liu/cypherscan-server/internal/util"
)

// BlockChain is the struct of the Client
type BlockChain struct {
	client  *ethclient.Client
	context context.Context
}

// BlockHeadersByNumbers is to get BlockHeaders by numbers
func (blockChain *BlockChain) BlockHeadersByNumbers(numbers []int64) ([]*types.Header, error) {
	return blockChain.client.BlockHeadersByNumbers(blockChain.context, numbers)
}

// BlockByHash is to get Block by hash, and return number of transactions without retreive whole transation slice
func (blockChain *BlockChain) BlockByHash(hash common.Hash, incTx bool) (*types.Block, int, error) {
	return blockChain.client.BlockByHash(blockChain.context, hash, incTx)
}

// KeyBlockByHash is to get Key Block by hash
func (blockChain *BlockChain) KeyBlockByHash(hash common.Hash) (*types.KeyBlock, error) {
	return blockChain.client.KeyBlockByHash(blockChain.context, hash)
}

// KeyBlocksByNumbers is to get BlockHeaders by numbers
func (blockChain *BlockChain) KeyBlocksByNumbers(numbers []int64) ([]*types.KeyBlock, error) {
	return blockChain.client.KeyBlocksByNumbers(blockChain.context, numbers)
}

// Subscribe is to subscirbe new block and new key block
func (blockChain *BlockChain) Subscribe(chBlock chan<- *types.Header, chKeyBlock chan<- *types.KeyBlockHeader) (Subscribed, error) {
	ctx := context.Background()
	blockSub, blockErr := blockChain.client.SubscribeNewHead(ctx, chBlock)
	if blockErr != nil {
		return nil, &util.MyError{Message: "Cannot subscirbe to Block"}
	}
	keyBlockSub, keyBlockErr := blockChain.client.SubscribeNewKeyHead(ctx, chKeyBlock)
	if keyBlockErr != nil {
		blockSub.Unsubscribe()
		return nil, &util.MyError{Message: "Cannot subscirbe to Key Block"}
	}

	return &CypherSubscribed{blockSub, keyBlockSub}, nil
}
