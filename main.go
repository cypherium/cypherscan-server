package main

import (
	"fmt"
	"gitlab.com/ron-liu/cypherscan-server/config"
)

func main() {
	fmt.Println("in main", configLib.Config.DbDrive)
}

