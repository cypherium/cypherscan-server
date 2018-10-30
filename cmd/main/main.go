package main

import (
  "fmt"

  "github.com/gin-gonic/gin"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "gitlab.com/ron-liu/cypherscan-server/internal/env"
  "gitlab.com/ron-liu/cypherscan-server/internal/home"
  "gitlab.com/ron-liu/cypherscan-server/internal/txblock"
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
)

func initDb() {
  db := util.OpenDb()
  db.AutoMigrate(&txblock.TxBlock{})
  defer db.Close()
}

func main() {
  initDb()

  fmt.Println("Evironments:", env.Env)
  routers := gin.Default()
  routers.GET("/home", home.GetHome)
  routers.Run()
}
