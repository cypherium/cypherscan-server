package repo

import (
	"fmt"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"sort"
	"strconv"
)

// QueryAddressRequest is the struct to pass the query options
type QueryAddressRequest struct {
	CursorPaginationRequest
	Address []byte
}

type SortTxByTime []*Transaction

func (s SortTxByTime) Len() int { return len(s) }
func (s SortTxByTime) Less(i, j int) bool {
	return s[i].Block.Time.Before(s[j].Block.Time)
}
func (s SortTxByTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// QueryAddress is to query transaction by send or recv address
func (r *Repo) QueryAddress(request *QueryAddressRequest) (*QueryResult, error) {
	var txs SortTxByTime
	var skip int64
	var cursor Cursor = Cursor{"", "1"}
	log.Info("QueryAddress", fmt.Sprintf("%s", request.CursorPaginationRequest))
	if request.CursorPaginationRequest.Before != "" {
		skip, _ = strconv.ParseInt(request.CursorPaginationRequest.Before, 10, 64)
		if skip > 0 {
			cursor.First = strconv.FormatInt(skip-1, 10)
			cursor.Last = strconv.FormatInt(skip+1, 10)
		}
	} else if request.CursorPaginationRequest.After != "" {
		skip, _ = strconv.ParseInt(request.CursorPaginationRequest.After, 10, 64)
		cursor.First = strconv.FormatInt(skip-1, 10)
		cursor.Last = strconv.FormatInt(skip+1, 10)
	}

	// err := r.dbRunner.Run(func(db *gorm.DB) error {
	// 	return db.Debug().Where("\"from\" = ?", request.Address).Or("\"to\" = ?", request.Address).Order("id desc").Offset(skip).Limit(request.CursorPaginationRequest.PageSize).Find(&txs).Error
	// })

	err := r.dbRunner.Run(func(db *gorm.DB) error {
		return db.Debug().Preload("Block", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"time", "number"})
		}).Where("\"from\" = ?", request.Address).Or("\"to\" = ?", request.Address).Order("id desc").Offset(skip * int64(request.CursorPaginationRequest.PageSize)).Limit(request.CursorPaginationRequest.PageSize).Find(&txs).Error
	})
	if err != nil {
		return nil, err
	}
	if request.CursorPaginationRequest.PageSize > len(txs) {
		if request.CursorPaginationRequest.Before != "" {
			cursor.First = ""
		} else {
			cursor.Last = ""
		}
	}
	sort.Sort(SortTxByTime(txs))
	log.Info("cursor", fmt.Sprintf("%s", cursor))
	log.Info("txs", fmt.Sprintf("%s", txs))
	log.Info("txs len", fmt.Sprintf("%s", len(txs)))
	return &QueryResult{
		Cursor: cursor,
		Items:  txs,
	}, nil

}
