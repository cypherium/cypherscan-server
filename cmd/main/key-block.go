package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	"github.com/cypherium/cypherscan-server/internal/repo"
	"github.com/gorilla/mux"
)

func getKeyBlocks(a *App, w http.ResponseWriter, r *http.Request) {
	// get request
	pageNo, pageSize, err := getPaginationRequest(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	latestNumber, err := a.blocksFetcher.GetLatestKeyBlockNumber()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var startWith = latestNumber - (pageNo-1)*int64(pageSize)
	keyBlocks, err := a.repo.GetKeyBlocks(&repo.BlockSearchContdition{Scenario: repo.ListPage, StartWith: startWith, PageSize: pageSize})
	dbListKeyBlocks := func(bs []repo.KeyBlock) []*listKeyBlock {
		ret := make([]*listKeyBlock, 0, len(keyBlocks))
		for _, b := range bs {
			ret = append(ret, &listKeyBlock{Number: b.Number, Time: b.Time, Difficulty: b.Difficulty})
		}
		return ret
	}(keyBlocks)
	numbersAlreadyGot := func() []int64 {
		ret := make([]int64, 0, len(keyBlocks))
		for _, b := range keyBlocks {
			ret = append(ret, b.Number)
		}
		return ret
	}()
	missedListKeyBlocks := func() []*listKeyBlock {
		if pageSize == len(numbersAlreadyGot) {
			return []*listKeyBlock{}
		}
		missedNumber := getMissedNumbers(latestNumber-int64(pageSize)*(pageNo-1), pageSize, numbersAlreadyGot)
		missedBlocks, _ := a.blocksFetcher.KeyBlocksByNumbers(missedNumber)
		return func(bs []*types.KeyBlock) []*listKeyBlock {
			ret := make([]*listKeyBlock, 0, len(keyBlocks))
			for _, h := range bs {
				ret = append(ret, transferKeyBlockHeadToListKeyBlock(h))
			}
			return ret

		}(missedBlocks)
	}()
	retList := append(dbListKeyBlocks, missedListKeyBlocks...)
	sort.Sort(numberDescSorterForListKeyBlock(retList))

	respondWithJSON(w, http.StatusOK, &responseOfGetKeyBlocks{Total: latestNumber + 1, Blocks: retList})
}

func getKeyBlock(a *App, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strNumber := vars["number"]
	number, err := strconv.ParseInt(strNumber, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("The passed number(%s) is not a valid number", strNumber))
		return
	}

	block, err := a.repo.GetKeyBlock(number)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, block)
}

type responseOfGetKeyBlocks struct {
	Total  int64           `json:"total"`
	Blocks []*listKeyBlock `json:"records"`
}

type listKeyBlock struct {
	Number     int64       `json:"number"`
	Time       time.Time   `json:"createdAt"`
	Difficulty repo.UInt64 `json:"difficulty"`
}

func transferKeyBlockHeadToListKeyBlock(h *types.KeyBlock) *listKeyBlock {
	return &listKeyBlock{
		Number:     h.Number().Int64(),
		Time:       time.Unix(h.Time().Int64(), 0),
		Difficulty: repo.UInt64(h.Difficulty().Uint64()),
	}
}

type numberDescSorterForListKeyBlock []*listKeyBlock

func (a numberDescSorterForListKeyBlock) Len() int           { return len(a) }
func (a numberDescSorterForListKeyBlock) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a numberDescSorterForListKeyBlock) Less(i, j int) bool { return a[i].Number > a[j].Number }
