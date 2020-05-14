package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cypherium/cypherscan-server/internal/repo"
	"github.com/cypherium/cypherscan-server/internal/util"
	"github.com/gorilla/mux"
)

func getOffsetPaginationRequest(r *http.Request, latestNo int64) (*repo.PaginationRequest, error) {
	pageSize, pageSizeErr := getPagesize(r)
	if pageSizeErr != nil {
		return nil, pageSizeErr
	}
	const (
		DefaultPageNo = "1"
	)
	v := r.URL.Query()
	strPageNo := v.Get("p")
	if strPageNo == "" {
		strPageNo = DefaultPageNo
	}
	pageNo, pageNoErr := strconv.ParseInt(strPageNo, 10, 64)
	if pageNoErr != nil {
		return nil, &util.MyError{Message: fmt.Sprintf("The passed pageNo(%s) is not a valid number", strPageNo)}
	}
	return &repo.PaginationRequest{
		PageSize: pageSize,
		PageNo:   pageNo,
		Type:     repo.OffsetType,
		LatestNo: latestNo,
	}, nil
}

func getCursorPaginationRequest(r *http.Request) (*repo.CursorPaginationRequest, error) {
	v := r.URL.Query()
	after := v.Get("after")
	before := v.Get("before")
	pageSize, pageSizeErr := getPagesize(r)
	if pageSizeErr != nil {
		return nil, pageSizeErr
	}
	return &repo.CursorPaginationRequest{After: after, PageSize: pageSize, Before: before}, nil
}

func getPagesize(r *http.Request) (int, error) {
	const (
		DefaultListPageSize = "20"
	)
	v := r.URL.Query()
	strPageSize := v.Get("pagesize")
	if strPageSize == "" {
		strPageSize = DefaultListPageSize
	}
	pageSize, pageSizeErr := strconv.Atoi(strPageSize)
	if pageSizeErr != nil {
		return 0, &util.MyError{Message: fmt.Sprintf("The passed pagesize(%s)  is not a valid number", strPageSize)}
	}
	return pageSize, nil
}

func getHashRequest(r *http.Request, name string) ([]byte, error) {
	const (
		DefaultListPageSize = "20"
	)
	vars := mux.Vars(r)
	strHash := vars[name]
	strBlockHash := fmt.Sprintf("\"%s\"", strHash)
	blockHash := Bytes(make([]byte, 100))
	blockHashErr := json.Unmarshal([]byte(strBlockHash), &blockHash)
	if blockHashErr != nil {
		return nil, &util.MyError{
			Inner:   blockHashErr,
			Message: fmt.Sprintf("The passed hash(%s) is not valid", strBlockHash),
		}
	}
	return blockHash, nil
}

func getNumberRequest(r *http.Request, name string) (int64, error) {
	vars := mux.Vars(r)
	strNumber := vars[name]
	number, err := strconv.ParseInt(strNumber, 10, 64)
	if err != nil {
		return 0, util.NewError(err, "The passed %s(%s) is not a number", name, strNumber)
	}
	return number, nil
}

// func respondWithError(w http.ResponseWriter, code int, message string) {
// 	respondWithJSON(w, code, map[string]string{"error": message})
// }

// func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

// 	response, err := json.Marshal(payload)
// 	if err != nil {
// 		log.Error(err.Error())
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Error when marshal object to json string"))
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	w.Write(response)
// }

// CursoredList is used to hold any list of data to frontend to display
type CursoredList struct {
	Items []interface{} `json:"records"`
	Last  string        `json:"last"`
	First string        `json:"first"`
}

// OffsetedList is sued to hold offseted list
type OffsetedList struct {
	Items      []interface{} `json:"records"`
	TotalCount int64         `json:"total"`
}
