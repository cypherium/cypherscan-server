package main

import (
  "fmt"
  // "github.com/gin-gonic/gin"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "gitlab.com/ron-liu/cypherscan-server/internal/env"
  "gitlab.com/ron-liu/cypherscan-server/internal/home"
  "gitlab.com/ron-liu/cypherscan-server/internal/txblock"
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
)

func initDb() {
  util.Run(func(db *gorm.DB) error {
    db.AutoMigrate(&txblock.TxBlock{}, &txblock.Transaction{})
    return nil
  })

}

func main() {
  fmt.Println("Evironments:", env.Env)
  util.OpenDb()
  initDb()
  defer util.CloseDb()
  home.SubscribeNewBlock()
  // routers := gin.Default()
  // routers.GET("/home", home.GetHome)
  // routers.Run()
}
