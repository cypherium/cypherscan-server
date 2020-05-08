package main

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/cypherium/cypherscan-syncer/internal/bizutil"
	"github.com/cypherium/cypherscan-syncer/internal/repo"
	"github.com/cypherium/cypherscan-syncer/internal/util"
)

func getSearch(a *App, w http.ResponseWriter, r *http.Request) {
	number, err := getNumberRequest(r, "q")
	if err == nil {
		block, keyBlock := a.repo.QueryNumber(number)
		respondWithJSON(w, 200, convertToNumberResult(block, keyBlock))
		return
	}

	bytes, err := getHashRequest(r, "q")
	if err != nil {
		bizutil.HandleError(err, "query failed")
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("GET /search/{q} failed"))
		return
	}
	bytesType := getBytesType(bytes)
	if bytesType == unknownType {
		bizutil.HandleError(err, "query failed")
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("GET /search/{q} failed"))
		return
	}

	if bytesType == addressType {
		pagination, err := getCursorPaginationRequest(r)
		if err != nil {
			bizutil.HandleError(err, "get cursor pagination request failed")
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("GET /search/{q} failed, because failed to get pagination queries"))
			return
		}
		queryResult, err := a.repo.QueryAddress(&repo.QueryAddressRequest{Address: bytes, CursorPaginationRequest: *pagination})
		if err != nil {
			bizutil.HandleError(err, "query address failed")
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("GET /search/{q} failed, because failed to get query address: %#v and pagination: %#v", bytes, pagination))
			return
		}
		client, err := a.pool.Borrow()
		if err != nil {
			bizutil.HandleError(util.NewError(err, "Borrow from pool failed"), "error when fetch block")
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("GET /search/%d failed, because failed to get balance", bytes))
			return
		}
		defer a.pool.Return(client)
		balance, err := client.Client.GetBalance(bytes)
		if err != nil {
			bizutil.HandleError(util.NewError(err, "GetBalance failed"), "")
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("GET /search/%d failed, because failed to get balance", bytes))
			return
		}
		respondWithJSON(w, 200, convertToAddressSearchResult(queryResult, balance))
		return
	}

	blockItem, err := a.repo.QueryHash(bytes)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("GET /search/{q} failed, because %#v", err))
		return
	}
	respondWithJSON(w, 200, convertToHashSearchResult(blockItem))
}

type bytesType int

const (
	addressType bytesType = 1
	hashType    bytesType = 2
	unknownType bytesType = 3
)

func getBytesType(bytes []byte) bytesType {
	if len(bytes) == 32 {
		return hashType
	}
	if len(bytes) == 20 {
		return addressType
	}
	return unknownType
}

type searchResultType string

const (
	number  searchResultType = "number"
	hash    searchResultType = "hash"
	address searchResultType = "address"
)

type searchResult struct {
	ResultType searchResultType `json:"type"`
	Result     interface{}      `json:"result"`
	Balance    string           `json:"balance"`
}

type numberResult struct {
	Block    *txBlock  `json:"block"`
	KeyBlock *keyBlock `json:"keyBlock"`
}

func convertToNumberResult(block *repo.BlockItem, keyBlock *repo.BlockItem) *searchResult {
	return &searchResult{
		ResultType: number,
		Result: &numberResult{
			Block:    convertToTxBlock(block),
			KeyBlock: convertToKeyBlock(keyBlock),
		},
	}
}
func convertToAddressSearchResult(queryResult *repo.QueryResult, balance *big.Int) *searchResult {
	return &searchResult{
		ResultType: address,
		Result:     convertQueryResultToListTxs(queryResult),
		Balance:    fmt.Sprintf("%d", balance),
	}
}

type hashSearchResultType string

const (
	blockHash    hashSearchResultType = "block"
	keyBlockHash hashSearchResultType = "keyBlock"
	txHash       hashSearchResultType = "tx"
)

type hashResult struct {
	HashResultType hashSearchResultType `json:"type"`
	Item           interface{}          `json:"item"`
}

func convertToHashSearchResult(b *repo.BlockItem) *searchResult {
	var result *hashResult
	if b == nil {
		return &searchResult{
			ResultType: hash,
			Result:     nil,
		}
	}
	if b.Type == repo.BlockType {
		result = &hashResult{
			HashResultType: blockHash,
			Item:           convertToTxBlock(b),
		}
	} else if b.Type == repo.KeyBlockType {
		result = &hashResult{
			HashResultType: keyBlockHash,
			Item:           convertToKeyBlock(b),
		}
	} else {
		result = &hashResult{
			HashResultType: txHash,
			Item:           convertToTx(b),
		}
	}
	return &searchResult{
		ResultType: hash,
		Result:     result,
	}
}
