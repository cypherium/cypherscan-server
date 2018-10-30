package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type _EnvStruct struct {
	DbDrive  string
	DbSource string
}

// Env is a structure holding env

var Env _EnvStruct

func init() {
	var err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Env = _EnvStruct{
		os.Getenv("DB_DRIVE"),
		os.Getenv("DB_SOURCE"),
	}
}
