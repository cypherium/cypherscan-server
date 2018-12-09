package util

import (
  "github.com/jinzhu/gorm"

  "gitlab.com/ron-liu/cypherscan-server/internal/env"
  "log"
)

var db *gorm.DB

// RunFunc is the type of function to run db scripts
type RunFunc func(db *gorm.DB) error

func getDb() *gorm.DB {
  return db
}

// OpenDb open db connection
func OpenDb() {
  _db, err := gorm.Open(env.Env.DbDrive, env.Env.DbSource)
  if err != nil {
    log.Fatalln("connect to db failed", err)
  }
  db = _db
}

// CloseDb close db connection
func CloseDb() {
  db.Close()
}

// RunDb a function take the db as the argument
func RunDb(fn RunFunc) error {
  db := getDb()
  return fn(db)
}
