package api

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
	"gitlab.com/ron-liu/cypherscan-server/internal/util"
)

// GetBlocks is get blocks based on to block number and the page size
func GetBlocks(w http.ResponseWriter, r *http.Request) {
	strNumber := mux.Vars(r)["number"]
	strPageSize := r.FormValue("pagesize")

	number, numberErr := util.StringToBigInt(strNumber, 10)
	pageSize, pageSizeErr := strconv.ParseInt(strPageSize, 10, 32)
	if numberErr != nil || pageSizeErr != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("The passed number or pageSize is not a valid number", strNumber))
		return
	}

	txBlocks := getTxBlocksFromDb(number, int(pageSize))
	respondWithJSON(w, http.StatusOK, txBlocks)
}

func getTxBlocksFromDb(number *big.Int, pageSize int) []repo.TxBlock {
	var txBlocks []repo.TxBlock
	max := number.Int64()
	min := max - int64(pageSize) + 1
	util.RunDb(func(db *gorm.DB) error {
		db.Where("number BETWEEN ? AND ?", min, max).Select([]string{"number", "txn", "time"}).Order("time desc").Limit(pageSize).Find(&txBlocks)
		return nil
	})
	return txBlocks
}
