package repo

import (
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
	var skip string
	if request.CursorPaginationRequest.Before != "" {
		skip = request.CursorPaginationRequest.Before
	} else if request.CursorPaginationRequest.After != "" {
		skip = request.CursorPaginationRequest.After
	} else {
		skip = "1"
	}

	err := r.dbRunner.Run(func(db *gorm.DB) error {
		return db.Debug().Where("\"from\" = ?", request.Address).Or("\"to\" = ?", request.Address).Order("id desc").Offset(skip).Limit(request.CursorPaginationRequest.PageSize).Find(&txs).Error
	})
	if err != nil {
		return nil, err
	}
	return &QueryResult{
		Items: txs,
	}, nil

}
