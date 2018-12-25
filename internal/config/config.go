package config

import (
	log "github.com/sirupsen/logrus"

	"os"

	"github.com/joho/godotenv"
)

// Config is the paramters
type Config struct {
	DbDrive         string
	DbSource        string
	BlockChainWsURL string
	OriginAllowed   string
}

// GetFromEnv is to get env
func GetFromEnv() *Config {
	var err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		os.Getenv("DB_DRIVE"),
		os.Getenv("DB_SOURCE"),
		os.Getenv("TX_BLOCKCHAIN_WS_URL"),
		os.Getenv("ORIGIN_ALLOWED"),
	}
}
