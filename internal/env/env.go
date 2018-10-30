package env

import (
  "log"
  "os"

  "github.com/joho/godotenv"
)

type _EnvStruct struct {
  DbDrive           string
  DbSource          string
  TsBlockChainWsURL string
}

// Env is a structure holding env
var Env _EnvStruct

func init() {
  var err = godotenv.Load()
  if err != nil {
    log.Println("Error loading .env file")
  }
  Env = _EnvStruct{
    os.Getenv("DB_DRIVE"),
    os.Getenv("DB_SOURCE"),
    os.Getenv("TX_BLOCKCHAIN_WS_URL"),
  }
}
