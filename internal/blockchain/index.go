package blockchain

import (
	"math/big"

	"github.com/cypherium/cypherBFT/common"
)

// Close is to close the client
func (blockChain *BlockChain) Close() {

}

// GetBalance is to get the latest account balance
func (blockChain *BlockChain) GetBalance(account []byte) (*big.Int, error) {
	var address common.Address
	copy(address[:], account[:])
	return blockChain.client.BalanceAt(blockChain.context, address, nil)
}
