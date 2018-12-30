package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ron-liu/cypherscan-server/internal/blockchain"
	"gitlab.com/ron-liu/cypherscan-server/internal/publisher"
	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
)

// App is the application structuer
type App struct {
	repo          repo.Get
	wsServer      publisher.WebSocketServer
	blocksFetcher blockchain.BlocksFetcher
	Router        *mux.Router
	originAllowed string
}

// NewApp is the constructor for App
func NewApp(rep repo.Get, wsServer publisher.WebSocketServer, blocksFetcher blockchain.BlocksFetcher, originAllowed string) *App {
	a := App{rep, wsServer, blocksFetcher, mux.NewRouter(), originAllowed}
	// a.setupCors()
	a.initializeRoutes()
	return &a
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/home", cors(a.GetHome)).Methods("GET", "OPTIONS")
	a.Router.HandleFunc("/ws", a.wsServer.ServeWebsocket)
	a.Router.Path("/tx-blocks/{number:[0-9]+}").Queries("pagesize", "{pagesize}").HandlerFunc(cors(a.GetBlocks)).Methods("GET", "OPTIONS")
}

// GetHome is: GET /home
func (a *App) GetHome(w http.ResponseWriter, r *http.Request) {
	txBlocks, err := a.repo.GetBlocks(&repo.BlockSearchContdition{Scenario: repo.HomePage})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	keyBlocks, err := a.repo.GetKeyBlocks(&repo.BlockSearchContdition{Scenario: repo.HomePage})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	transactions, err := a.repo.GetTransactions(&repo.TransactionSearchCondition{Scenario: repo.HomePage})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	payload := HomePayload{
		Metrics: []HomeMetric{},
		TxBlocks: func() []HomeTxBlock {
			ret := make([]HomeTxBlock, 0, len(txBlocks))
			for _, b := range txBlocks {
				ret = append(ret, HomeTxBlock{b.Number, b.Txn, b.Time})
			}
			return ret
		}(),
		KeyBlocks: func() []HomeKeyBlock {
			ret := make([]HomeKeyBlock, 0, len(keyBlocks))
			for _, b := range keyBlocks {
				ret = append(ret, HomeKeyBlock{b.Number, b.Time})
			}
			return ret
		}(),
		Txs: func() []HomeTx {
			ret := make([]HomeTx, 0, len(transactions))
			for _, t := range transactions {
				ret = append(ret, HomeTx{
					t.Block.Time,
					t.Value,
					t.Hash.Hex(),
					t.From.Hex(),
					t.To.Hex(),
				})
			}
			return ret
		}(),
	}
	respondWithJSON(w, http.StatusOK, payload)
}

// GetBlocks is: GET /tx-blocks/:{number}?pagesize={pageszie}
func (a *App) GetBlocks(w http.ResponseWriter, r *http.Request) {
	strNumber := mux.Vars(r)["number"]
	strPageSize := r.FormValue("pagesize")

	number, numberErr := strconv.ParseInt(strNumber, 10, 64)
	pageSize, pageSizeErr := strconv.Atoi(strPageSize)
	if numberErr != nil || pageSizeErr != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("The passed number or pageSize is not a valid number", strNumber))
		return
	}
	txBlocks, _ := a.repo.GetBlocks(&repo.BlockSearchContdition{Scenario: repo.ListPage, StartWith: number, PageSize: pageSize})
	dbListTxBlocks := func(bs []repo.TxBlock) []*listTxBlock {
		ret := make([]*listTxBlock, 0, len(txBlocks))
		for _, b := range bs {
			ret = append(ret, &listTxBlock{Number: b.Number, Time: b.Time, GasUsed: uint64(b.GasUsed), GasLimit: uint64(b.GasLimit)})
		}
		return ret
	}(txBlocks)
	numbersAlreadyGot := func() []int64 {
		ret := make([]int64, 0, len(txBlocks))
		for _, b := range txBlocks {
			ret = append(ret, b.Number)
		}
		return ret
	}()

	missedListTxBlocks := func() []*listTxBlock {
		if pageSize == len(numbersAlreadyGot) {
			return []*listTxBlock{}
		}
		missedNumber := getMissedNumbers(number, pageSize, numbersAlreadyGot)
		missedBlocks, _ := a.blocksFetcher.BlockHeadersByNumbers(missedNumber)
		return func(bs []*types.Header) []*listTxBlock {
			ret := make([]*listTxBlock, 0, len(txBlocks))
			for _, h := range bs {
				ret = append(ret, transferBlockHeadToListTxBlock(h))
			}
			return ret

		}(missedBlocks)
	}()
	retList := append(dbListTxBlocks, missedListTxBlocks...)
	sort.Sort(numberDescSorterForListTxBlock(retList))

	respondWithJSON(w, http.StatusOK, retList)
}

// Run starts the app and serves on the specified addr
func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}
