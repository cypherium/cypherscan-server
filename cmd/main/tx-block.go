package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	"github.com/gorilla/mux"
	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
	"gitlab.com/ron-liu/cypherscan-server/internal/util"
)

func getBlock(a *App, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strNumber := vars["number"]
	number, err := strconv.ParseInt(strNumber, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("The passed number(%s) is not a valid number", strNumber))
		return
	}

	block, err := a.repo.GetBlock(number)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, block)
}

func getBlocks(a *App, w http.ResponseWriter, r *http.Request) {
	// get request
	pageNo, pageSize, err := func(r *http.Request) (int64, int, error) {
		const (
			DefaultPageNo       = "1"
			DefaultListPageSize = "20"
		)
		v := r.URL.Query()
		strPageNo := v.Get("p")
		strPageSize := v.Get("pagesize")
		if strPageNo == "" {
			strPageNo = DefaultPageNo
		}
		if strPageSize == "" {
			strPageSize = DefaultListPageSize
		}

		pageNo, pageNoErr := strconv.ParseInt(strPageNo, 10, 64)
		pageSize, pageSizeErr := strconv.Atoi(strPageSize)
		if pageNoErr != nil || pageSizeErr != nil {
			return 0, 0, &util.MyError{Message: fmt.Sprintf("The passed p(%s) or pageSize(%s) is not a valid number", strPageNo, strPageSize)}
		}
		return pageNo, pageSize, nil
	}(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	latestNumber, err := a.blocksFetcher.GetLatestBlockNumber()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var startWith = latestNumber - (pageNo-1)*int64(pageSize)
	txBlocks, err := a.repo.GetBlocks(&repo.BlockSearchContdition{Scenario: repo.ListPage, StartWith: startWith, PageSize: pageSize})
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
		missedNumber := getMissedNumbers(latestNumber-int64(pageSize)*(pageNo-1), pageSize, numbersAlreadyGot)
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

	respondWithJSON(w, http.StatusOK, &ResponseOfGetBlocks{Total: latestNumber + 1, Blocks: retList})

}

type ResponseOfGetBlocks struct {
	Total  int64          `json:"total"`
	Blocks []*listTxBlock `json:"blocks"`
}

type listTxBlock struct {
	Number   int64     `json:"number"`
	Time     time.Time `json:"createdAt"`
	Txn      int       `json:"txn"`
	GasUsed  uint64    `json:"gasUsed"`
	GasLimit uint64    `json:"gasLimit"`
}

func transferBlockHeadToListTxBlock(h *types.Header) *listTxBlock {
	return &listTxBlock{
		Number: h.Number.Int64(),
		Time:   time.Unix(h.Time.Int64(), 0),
		// txn
		GasUsed:  h.GasUsed,
		GasLimit: h.GasLimit,
	}
}

type numberDescSorterForListTxBlock []*listTxBlock

func (a numberDescSorterForListTxBlock) Len() int           { return len(a) }
func (a numberDescSorterForListTxBlock) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a numberDescSorterForListTxBlock) Less(i, j int) bool { return a[i].Number > a[j].Number }
