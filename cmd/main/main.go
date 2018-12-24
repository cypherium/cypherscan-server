package main

import (
	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ron-liu/cypherscan-server/internal/api"
	"gitlab.com/ron-liu/cypherscan-server/internal/env"
	"gitlab.com/ron-liu/cypherscan-server/internal/publisher"
	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
	"gitlab.com/ron-liu/cypherscan-server/internal/util"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Evironments:", env.Env)
	// context := context.Background()

	dbClient, err := util.Connect()
	if err != nil {
		log.Fatal("Can NOT connect to database")
	}
	defer dbClient.Close()

	repo := repo.NewRepo(dbClient)
	repo.InitDb()

	// blockChainClient, err := blockchain.Dial(context)
	// if err != nil {
	// 	log.Fatal("Can NOT connect to blockchain")
	// }

	// go blockchain.SubscribeNewBlock([]blockchain.BlockHandlers{repo.SaveBlock, boardcastNewBlock})
	// go publisher.StartHub()

	controller := api.NewController(repo)
	util.CreateRouter(func(r *mux.Router) {
		r.HandleFunc("/home", controller.GetHome).Methods("GET")
		r.HandleFunc("/ws", api.HanderForBrowser)
		r.Path("/tx-blocks/{number:[0-9]+}").Queries("pagesize", "{pagesize}").HandlerFunc(controller.GetBlocks)
	})
}

func boardcastNewBlock(block *types.Block) error {
	publisher.Broadcast(api.TransformTxBlockToFrontendMessage(block))
	return nil
}
