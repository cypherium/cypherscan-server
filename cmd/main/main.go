package main

import (
	"context"
	"fmt"

	"github.com/cypherium/cypherBFT/go-cypherium/core/types"
	"github.com/cypherium/cypherscan-server/internal/blockchain"
	"github.com/cypherium/cypherscan-server/internal/config"
	"github.com/cypherium/cypherscan-server/internal/publisher"
	"github.com/cypherium/cypherscan-server/internal/repo"
	"github.com/cypherium/cypherscan-server/internal/util"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)
import (
	"github.com/jet/go-interstellar"
	"net/http"
	"os"
	"testing"
)

// CreateTestClient creates an *interstellar.Client for tests
// It gets the cosmos db connection string from the environment variable `AZURE_COSMOS_DB_CONNECTION_STRING`
// If the environment variable is not set, it will cause the test to be skipped.
// If the environment variable fails to parse, the test will fail.
func CreateTestClient(t *testing.T) *interstellar.Client {
	cstring := os.Getenv("AZURE_COSMOS_DB_CONNECTION_STRING")

	if cstring == "" {
		t.Skip("Must provide AZURE_COSMOS_DB_CONNECTION_STRING environment variable to test")
	}

	cs, err := interstellar.ParseConnectionString(cstring)
	if err != nil {
		t.Fatal(err)
	}
	client, _ := interstellar.NewClient(cs, NewTestLoggingRequester(t, http.DefaultClient))
	return client
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	context := context.Background()
	config := config.GetFromEnv()
	log.Info("Config:", fmt.Sprintf("%v", config))

	dbClient, err := util.ConnectDb("sqlite3", config.RdsHostName, config.RdsPort, config.RdsDbName, config.RdsUserName, config.RdsPassword, config.RdsSslMode)
	if err != nil {
		log.Fatal(fmt.Sprintf("Can NOT connect to database: %s", err.Error()))
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

	app := NewApp(repoInstance, hub, blockChainClient, config.OriginAllowed)
	app.Run()
}
