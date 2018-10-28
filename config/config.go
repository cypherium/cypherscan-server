package configLib

import (
	"github.com/joho/godotenv"
	"os"
	"log"
)

type ConfigStruct struct {
	DbDrive string
	DbSource string
}
var Config ConfigStruct 

func init () {
	var err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Config = ConfigStruct {
		os.Getenv("DB_DRIVE"),
		os.Getenv("DB_SOURCE"),
	}
}
