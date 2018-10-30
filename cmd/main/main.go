package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/ron-liu/cypherscan-server/internal/env"
	"gitlab.com/ron-liu/cypherscan-server/internal/home"
)

func main() {
	fmt.Println("Evironments:", env.Env)
	routers := gin.Default()
	routers.GET("/home", home.GetHome)
	routers.Run()
}
