package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ron-liu/cypherscan-server/internal/publisher"
	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
)

// App is the application structuer
type App struct {
	repo          repo.Get
	wsServer      publisher.WebSocketServer
	Router        *mux.Router
	originAllowed string
}

// NewApp is the constructor for App
func NewApp(rep repo.Get, wsServer publisher.WebSocketServer, originAllowed string) *App {
	a := App{rep, wsServer, mux.NewRouter(), originAllowed}
	a.initializeRoutes()
	return &a
}

func (a *App) initializeRoutes() {
	// Handle all preflight request
	a.Router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.WriteHeader(http.StatusNoContent)
		return
	})
	a.Router.StrictSlash(true)
	_ = handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	_ = handlers.AllowedOrigins([]string{a.originAllowed})
	_ = handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	a.Router.HandleFunc("/home", a.GetHome).Methods("GET")
	a.Router.HandleFunc("/ws", a.wsServer.ServeWebsocket)
	a.Router.Path("/tx-blocks/{number:[0-9]+}").Queries("pagesize", "{pagesize}").HandlerFunc(a.GetBlocks)
}

// GetHome is: GET /home
func (a *App) GetHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("starting getting home")
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
	payload := Payload{
		Metrics: []Metric{},
		TxBlocks: func() []TxBlock {
			ret := make([]TxBlock, 0, len(txBlocks))
			for _, b := range txBlocks {
				ret = append(ret, TxBlock{b.Number, b.Txn, b.Time})
			}
			return ret
		}(),
		KeyBlocks: func() []KeyBlock {
			ret := make([]KeyBlock, 0, len(keyBlocks))
			for _, b := range keyBlocks {
				ret = append(ret, KeyBlock{b.Number, b.Time})
			}
			return ret
		}(),
		Txs: func() []Tx {
			ret := make([]Tx, 0, len(transactions))
			for _, t := range transactions {
				ret = append(ret, Tx{
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
	respondWithJSON(w, http.StatusOK, txBlocks)
}

// Run starts the app and serves on the specified addr
func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

type setupRoute func(*mux.Router)

// CreateRouter is to generate the router
func CreateRouter(f setupRoute, originAllowed string) {
	r := mux.NewRouter()
	// Handle all preflight request
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.WriteHeader(http.StatusNoContent)
		return
	})
	r.StrictSlash(true)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{originAllowed})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	f(r)
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(methodsOk, headersOk, originsOk)(r)))
}
