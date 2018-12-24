package api

import (
	"fmt"
	"net/http"

	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
)

// Controller is a struct to hold all the routes' handlers
type Controller struct {
	repo repo.Get
}

// NewController is the constructor to create a controller
func NewController(repo repo.Get) *Controller {
	return &Controller{repo: repo}
}

// GetHome is: GET /home
func (c *Controller) GetHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("starting getting home")
	txBlocks, err := c.repo.GetBlocks(&repo.BlockSearchContdition{Scenario: repo.HomePage})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	keyBlocks, err := c.repo.GetKeyBlocks(&repo.BlockSearchContdition{Scenario: repo.HomePage})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	transactions, err := c.repo.GetTransactions(&repo.TransactionSearchCondition{Scenario: repo.HomePage})
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
func (c *Controller) GetBlocks(w http.ResponseWriter, r *http.Request) {
	strNumber := mux.Vars(r)["number"]
	strPageSize := r.FormValue("pagesize")

	number, numberErr := strconv.ParseInt(strNumber, 10, 64)
	pageSize, pageSizeErr := strconv.Atoi(strPageSize)
	if numberErr != nil || pageSizeErr != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("The passed number or pageSize is not a valid number", strNumber))
		return
	}

	txBlocks, _ := c.repo.GetBlocks(&repo.BlockSearchContdition{Scenario: repo.ListPage, StartWith: number, PageSize: pageSize})
	respondWithJSON(w, http.StatusOK, txBlocks)
}
