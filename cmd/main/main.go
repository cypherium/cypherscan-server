package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cypherium/cypherBFT/core/types"
	"github.com/cypherium/cypherscan-server/internal/bizutil"
	"github.com/cypherium/cypherscan-server/internal/blockchain"
	"github.com/cypherium/cypherscan-server/internal/config"

	"github.com/cypherium/cypherscan-server/internal/publisher"
	"github.com/cypherium/cypherscan-server/internal/repo"
	"github.com/cypherium/cypherscan-server/internal/util"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	context := context.Background()
	config := config.GetFromEnv()
	log.Info("Config:", fmt.Sprintf("%v", config))

	dbClient, err := util.ConnectDb("sqlite3", config.RdsHostName, config.RdsPort, "cypherdb", "postgres", "postgres", "disable")
	if err != nil {
		log.Fatal(fmt.Sprintf("Can NOT connect to database: %s", err.Error()))
	}
	defer dbClient.Close()

	executionTimeout, _ := strconv.Atoi(config.ExecutionTimeout)
	ctx, cancel := bizutil.GetContext(executionTimeout)

	repoInstance := repo.NewRepo(dbClient)
	repoInstance.InitDb()

	blockChainClient, err := blockchain.Dial(context, config.BlockChainWsURL)
	if err != nil {
		log.Fatal("Can NOT connect to blockchain")
		log.Info("err:", fmt.Sprintf("%v", err))
	}

	hub := publisher.NewHub()
	go hub.StartHub()

	chBlock := make(chan *types.Header)
	chKeyBlock := make(chan *types.KeyBlockHeader)
	newBlockListener := NewBlockListener{Repo: repoInstance, BlockFetcher: blockChainClient, Broadcastable: hub}
	go newBlockListener.Listen(chBlock, chKeyBlock)
	_, err = blockChainClient.Subscribe(chBlock, chKeyBlock)
	if err != nil {
		log.Fatal("Cannot subscribe blockchain")
	}
	pool, poolTerminaterd := blockchain.NewPool(ctx, &blockchain.NewPoolOptions{BorrowTimeoutMs: 5000, MaxSize: 2, NodesUrls: []string{config.BlockChainWsURL}})
	defer func() {
		cancel()
		<-poolTerminaterd
	}()
	app := NewApp(repoInstance, hub, blockChainClient, config.OriginAllowed, pool)
	app.Run()

}
