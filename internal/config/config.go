package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config is the paramters
type Config struct {
	DbDrive          string
	RdsHostName      string
	RdsPort          string
	RdsDbName        string
	RdsUserName      string
	RdsPassword      string
	RdsSslMode       string
	BlockChainWsURL  string
	OriginAllowed    string
	ExecutionTimeout string
}

// GetFromEnv is to get env
func GetFromEnv() *Config {
	godotenv.Load()
	return &Config{
		os.Getenv("DB_DRIVE"),
		"13.72.80.40",
		//"127.0.0.1",
		//os.Getenv("RDS_HOSTNAME"),
		//os.Getenv("RDS_PORT"),
		"8546",
		os.Getenv("RDS_DB_NAME"),
		os.Getenv("RDS_USERNAME"),
		os.Getenv("RDS_PASSWORD"),
		os.Getenv("RDS_SSLMODE"),
		//os.Getenv("TX_BLOCKCHAIN_WS_URL"),
		"ws://13.72.80.40:8546",
		//"ws://127.0.0.1:8546",
		os.Getenv("ORIGIN_ALLOWED"),
		"0",
	}
}
