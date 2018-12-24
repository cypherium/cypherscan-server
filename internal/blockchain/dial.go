package blockchain

import (
	"context"

	"github.com/cypherium/CypherTestNet/go-cypherium/ethclient"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ron-liu/cypherscan-server/internal/env"
)

// Dial is to connect the block chain
func Dial(context context.Context) (*BlockChain, error) {
	c, err := ethclient.Dial(env.Env.TsBlockChainWsURL)
	if err != nil {
		return nil, err
	}
	log.Info("Connected to blockchain nodes")
	return &BlockChain{c, context}, nil
}
