package util

import (
  "github.com/jinzhu/gorm"

  "gitlab.com/ron-liu/cypherscan-server/internal/env"
  "log"
)

var db *gorm.DB

// OpenDb open db connection
func OpenDb() {
  _db, err := gorm.Open(env.Env.DbDrive, env.Env.DbSource)
  if err != nil {
    log.Fatalln("connect to db failed")
  }
  db = _db
}

// GetDb return db connection
func GetDb() *gorm.DB {
  return db
}

// CloseDb close db connection
func CloseDb() {
  db.Close()
}
