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
	var skip int64
	var cursor repo.Cursor = repo.Cursor{"", "1"}

	pagination, err := getCursorPaginationRequest(r)
	if pagination.Before != "" {
		skip, _ = strconv.ParseInt(pagination.Before, 10, 64)
		if skip > 0 {
			cursor.First = strconv.FormatInt(skip-1, 10)
			cursor.Last = strconv.FormatInt(skip+1, 10)
		}
	} else if pagination.After != "" {
		skip, _ = strconv.ParseInt(pagination.After, 10, 64)
		cursor.First = strconv.FormatInt(skip-1, 10)
		cursor.Last = strconv.FormatInt(skip+1, 10)
	}
	txs, err := a.repo.GetTransactions(&repo.TransactionSearchCondition{BlockNumber: -1, PageSize: pagination.PageSize, Skip: skip * int64(pagination.PageSize), Scenario: repo.ListPage})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if pagination.PageSize > len(txs) {
		if pagination.Before != "" {
			cursor.First = ""
		} else {
			cursor.Last = ""
		}
	}

	respondWithJSON(w, http.StatusOK, convertQueryResultToListTxs(&repo.QueryResult{
		Cursor: cursor,
		Items:  txs}))
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
	respondWithJSON(w, http.StatusOK, convertToTx(tx))
}

type listTx struct {
	Time   time.Time    `json:"createdAt"`
	Value  uint64       `json:"value"`
	Hash   repo.Hash    `json:"hash"`
	From   repo.Address `json:"from"`
	To     repo.Address `json:"to"`
	Number int64        `json:"number"`
}

func transferTransactionToListTx(tx repo.Transaction) *listTx {
	return &listTx{
		Hash:   tx.Hash,
		Value:  uint64(tx.Value),
		From:   tx.From,
		To:     tx.To,
		Time:   tx.Block.Time,
		Number: tx.Block.Number,
	}
}

type tx struct {
	listTx
	GasPrice    uint64 `json:"gasPrice"`
	Gas         uint64 `json:"gas"`
	Cost        uint64 `json:"cost"`
	Payload     Bytes  `json:"input"`
	TxIndex     int    `json:"transactionIndex"`
	BlockHash   Bytes  `json:"blockHash"`
	BlockNumber int64  `json:"blockNumber"`
	Signature   Bytes  `json:"signature" `
}

func convertToTx(b *repo.Transaction) *tx {
	return &tx{
		listTx:      *transferTransactionToListTx(*b),
		GasPrice:    uint64(b.GasPrice),
		Gas:         uint64(b.Gas),
		Cost:        uint64(b.Cost),
		Payload:     b.Payload,
		TxIndex:     int(b.TransactionIndex),
		BlockHash:   b.BlockHash.Bytes(),
		BlockNumber: b.BlockNumber,
		Signature:   Bytes(b.Signature),
	}
}

type responseOfGetTxs struct {
	Total int64     `json:"total"`
	Txs   []*listTx `json:"records"`
}
