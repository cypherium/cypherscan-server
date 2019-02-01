package repo_test

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
	"gitlab.com/ron-liu/cypherscan-server/internal/util"
)

type testFunc func(db *gorm.DB)

func testOnAnCleanDb(f testFunc) {
	db, err := util.ConnectDb("postgres", "localhost", "5432", "scan_test", "postgres", "postgres", "disable")
	if err != nil {
		log.Fatal(fmt.Sprintf("Can NOT connect to database: %s", err.Error()))
	}
	db.Run(func(db *gorm.DB) error {
		db.DropTableIfExists(&repo.TxBlock{}, &repo.KeyBlock{}, &repo.Transaction{})
		db.AutoMigrate(&repo.TxBlock{}, &repo.Transaction{})
		return nil
	})
	defer db.Close()
	db.Run(func(db *gorm.DB) error {
		f(db)
		return nil
	})
}
