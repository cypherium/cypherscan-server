package main

import (
	"context"
	"fmt"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ron-liu/cypherscan-server/internal/blockchain"
	"gitlab.com/ron-liu/cypherscan-server/internal/config"
	"gitlab.com/ron-liu/cypherscan-server/internal/publisher"
	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
	"gitlab.com/ron-liu/cypherscan-server/internal/util"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	context := context.Background()
	config := config.GetFromEnv()
	log.Info("Config:", fmt.Sprintf("%v", config))

	dbClient, err := util.ConnectDb(config.DbDrive, config.DbSource)
	if err != nil {
		log.Fatal("Can NOT connect to database")
	}
	defer dbClient.Close()

	repoInstance := repo.NewRepo(dbClient)
	repoInstance.InitDb()

	blockChainClient, err := blockchain.Dial(context, config.BlockChainWsURL)
	if err != nil {
		log.Fatal("Can NOT connect to blockchain")
	}

	hub := publisher.NewHub()
	go hub.StartHub()

	newBlockListener := NewBlockListener{Repo: repoInstance, BlockFetcher: blockChainClient, Broadcastable: hub}
	chBlock := make(chan *types.Header)
	chKeyBlock := make(chan *types.KeyBlockHeader)
	_, err = blockChainClient.Subscribe(chBlock, chKeyBlock)
	if err != nil {
		log.Fatal("Cannot subscribe blockchain")
	}
	go newBlockListener.Listen(chBlock, chKeyBlock)

	app := NewApp(repoInstance, hub, config.OriginAllowed)
	app.Run()

}
