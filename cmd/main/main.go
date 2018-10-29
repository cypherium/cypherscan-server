package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/ron-liu/cypherscan-server/internal/config"
	"gitlab.com/ron-liu/cypherscan-server/internal/home"
)

func main() {
	fmt.Println("Evironments:", configLib.Config)
	routers := gin.Default()
	routers.GET("/home", home.GetHome)
	routers.Run()
}
