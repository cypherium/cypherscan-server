package repo

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

// QueryAddressRequest is the struct to pass the query options
type QueryAddressRequest struct {
	CursorPaginationRequest
	Address []byte
}

// QueryAddress is to query transaction by send or recv address
func (r *Repo) QueryAddress(request *QueryAddressRequest) (*QueryResult, error) {
	var txs []Transaction
	var skip int64
	var cursor Cursor = Cursor{"", "1"}

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
	return &QueryResult{
		Cursor: cursor,
		Items:  txs,
	}, nil

}
