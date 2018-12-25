package main

import (
	"context"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ron-liu/cypherscan-server/internal/api"
	"gitlab.com/ron-liu/cypherscan-server/internal/blockchain"
	"gitlab.com/ron-liu/cypherscan-server/internal/env"
	"gitlab.com/ron-liu/cypherscan-server/internal/publisher"
	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
	"gitlab.com/ron-liu/cypherscan-server/internal/util"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Evironments:", env.Env)
	context := context.Background()

	dbClient, err := util.Connect()
	if err != nil {
		log.Fatal("Can NOT connect to database")
	}
	defer dbClient.Close()

	repoInstance := repo.NewRepo(dbClient)
	repoInstance.InitDb()

	blockChainClient, err := blockchain.Dial(context)
	if err != nil {
		log.Fatal("Can NOT connect to blockchain")
	}

	hub := publisher.NewHub()
	go hub.StartHub()

	newBlockListener := api.NewBlockListener{Repo: repoInstance, BlockFetcher: blockChainClient, Broadcastable: hub}
	chBlock := make(chan *types.Header)
	chKeyBlock := make(chan *types.KeyBlockHeader)
	_, err = blockChainClient.Subscribe(chBlock, chKeyBlock)
	if err != nil {
		log.Fatal("Cannot subscribe blockchain")
	}
	go newBlockListener.Listen(chBlock, chKeyBlock)

	controller := api.NewController(repoInstance)
	util.CreateRouter(func(r *mux.Router) {
		r.HandleFunc("/home", controller.GetHome).Methods("GET")
		r.HandleFunc("/ws", hub.ServeWebsocket)
		r.Path("/tx-blocks/{number:[0-9]+}").Queries("pagesize", "{pagesize}").HandlerFunc(controller.GetBlocks)
	})
}
