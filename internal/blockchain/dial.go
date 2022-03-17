package blockchain

import (
	"context"

	"github.com/cypherium/cypherBFT/ethclient"

	log "github.com/sirupsen/logrus"
)

// Dial is to connect the block chain
func Dial(context context.Context, url string) (*BlockChain, error) {
	c, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	log.Info("Connected to blockchain nodes")
	return &BlockChain{c, context, -1, -1, -1, -1}, nil
}
