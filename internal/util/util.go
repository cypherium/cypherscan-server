package util

import (
  "github.com/jinzhu/gorm"

  "gitlab.com/ron-liu/cypherscan-server/internal/env"
  "log"
)

// OpenDb open db connection
func OpenDb() *gorm.DB {
  db, err := gorm.Open(env.Env.DbDrive, env.Env.DbSource)
  if err != nil {
    log.Fatalln("connect to db failed")
  }
  return db
}
