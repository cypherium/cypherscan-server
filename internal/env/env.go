package env

import (
  "log"
  "os"

  "github.com/joho/godotenv"
)

type envStruct struct {
  DbDrive           string
  DbSource          string
  TsBlockChainWsURL string
  OriginAllowed     string
}

// Env is a structure holding env
var Env envStruct

func init() {
  var err = godotenv.Load()
  if err != nil {
    log.Println("Error loading .env file")
  }
  Env = envStruct{
    os.Getenv("DB_DRIVE"),
    os.Getenv("DB_SOURCE"),
    os.Getenv("TX_BLOCKCHAIN_WS_URL"),
    os.Getenv("ORIGIN_ALLOWED"),
  }
}
