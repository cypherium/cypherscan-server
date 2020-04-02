package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cypherium/cypherscan-server/internal/repo"

	"github.com/gorilla/mux"
)

const (
	TotalTxsNumber = 1000
)

func getTxs(a *App, w http.ResponseWriter, r *http.Request) {
	pageNo, pageSize, err := getPaginationRequest(r)
	txs, err := a.repo.GetTransactions(&repo.TransactionSearchCondition{BlockNumber: -1, PageSize: pageSize, Skip: (pageNo - 1) * int64(pageSize), Scenario: repo.ListPage})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	list := make([]*listTx, 0, len(txs))
	for _, t := range txs {
		list = append(list, transferTransactionToListTx(t))
	}
	respondWithJSON(w, http.StatusOK, &responseOfGetTxs{Total: TotalTxsNumber, Txs: list})
}

func getBlockTxs(a *App, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strNumber := vars["number"]
	number, err := strconv.ParseInt(strNumber, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("The passed number(%s) is not a valid number", strNumber))
		return
	}
	pageNo, pageSize, err := getPaginationRequest(r)

	block, err := a.repo.GetBlock(number)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	txs, err := a.repo.GetTransactions(&repo.TransactionSearchCondition{BlockNumber: number, PageSize: pageSize, Skip: (pageNo - 1) * int64(pageSize), Scenario: repo.ListPage})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	list := make([]*listTx, 0, len(txs))
	for _, t := range txs {
		list = append(list, transferTransactionToListTx(t))
	}
	respondWithJSON(w, http.StatusOK, &responseOfGetTxs{Total: int64(block.Txn), Txs: list})
}

func getTx(a *App, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strHash := fmt.Sprintf("\"%s\"", vars["hash"])
	var hash repo.Hash
	err := json.Unmarshal([]byte(strHash), &hash)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("The passed hash(%s) is not valid", strHash))
		return
	}

	tx, err := a.repo.GetTransaction(hash)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, tx)
}

type listTx struct {
	CreatedAt time.Time    `json:"createdAt"`
	Value     uint64       `json:"value"`
	Hash      repo.Hash    `json:"hash"`
	From      repo.Address `json:"from"`
	To        repo.Address `json:"to"`
}

func transferTransactionToListTx(tx repo.Transaction) *listTx {
	return &listTx{
		Hash:      tx.Hash,
		Value:     uint64(tx.Value),
		From:      tx.From,
		To:        tx.To,
		CreatedAt: tx.Block.Time,
	}
}

type responseOfGetTxs struct {
	Total int64     `json:"total"`
	Txs   []*listTx `json:"records"`
}
