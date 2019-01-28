package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config is the paramters
type Config struct {
	DbDrive         string
	RdsHostName     string
	RdsPort         string
	RdsDbName       string
	RdsUserName     string
	RdsPassword     string
	RdsSslMode      string
	BlockChainWsURL string
	OriginAllowed   string
}

// GetFromEnv is to get env
func GetFromEnv() *Config {
	godotenv.Load()
	return &Config{
		os.Getenv("DB_DRIVE"),
		os.Getenv("RDS_HOSTNAME"),
		os.Getenv("RDS_PORT"),
		os.Getenv("RDS_DB_NAME"),
		os.Getenv("RDS_USERNAME"),
		os.Getenv("RDS_PASSWORD"),
		os.Getenv("RDS_SSLMODE"),
		os.Getenv("TX_BLOCKCHAIN_WS_URL"),
		os.Getenv("ORIGIN_ALLOWED"),
	}
}
