package main

import (
	"net/http"

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
	a.Router.Path("/tx-blocks").Queries("p", "{p}", "pagesize", "{pageSize}").HandlerFunc(cors(a.GetBlocks)).Methods("GET", "OPTIONS")
	a.Router.Path("/tx-blocks").HandlerFunc(cors(a.GetBlocks)).Methods("GET", "OPTIONS")
	a.Router.Path("/tx-block/{number:[0-9]+}").HandlerFunc(cors(a.GetBlock)).Methods("GET", "OPTIONS")
}

// GetHome is: GET /home
func (a *App) GetHome(w http.ResponseWriter, r *http.Request) {
	getHome(a, w, r)
}

// GetBlocks is: GET /tx-blocks/:{number}?pagesize={pageszie}
func (a *App) GetBlocks(w http.ResponseWriter, r *http.Request) {
	getBlocks(a, w, r)
}

// GetBlock is : GET /tx-block/{number}
func (a *App) GetBlock(w http.ResponseWriter, r *http.Request) {
	getBlock(a, w, r)
}

// GetKeyBlocks is: GET /tx-blocks/:{number}?pagesize={pageszie}
func (a *App) GetKeyBlocks(w http.ResponseWriter, r *http.Request) {
	getKeyBlocks(a, w, r)
}

// GetKeyBlock is : GET /key-block/{number}
func (a *App) GetKeyBlock(w http.ResponseWriter, r *http.Request) {
	getKeyBlock(a, w, r)
}

// Run starts the app and serves on the specified addr
func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}