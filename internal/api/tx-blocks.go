package api

import (
  "fmt"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  "gitlab.com/ron-liu/cypherscan-server/internal/txblock"
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
  "math/big"
  "net/http"
  "strconv"
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

func getTxBlocksFromDb(number *big.Int, pageSize int) []txblock.TxBlock {
  var txBlocks []txblock.TxBlock
  util.RunDb(func(db *gorm.DB) error {
    db.Select([]string{"number", "txn", "time"}).Order("time desc").Limit(TxBlockCount).Find(&txBlocks)
    return nil
  })
  return txBlocks
}
