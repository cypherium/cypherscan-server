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
		"119.8.146.229",
		"8546",
		os.Getenv("RDS_DB_NAME"),
		os.Getenv("RDS_USERNAME"),
		os.Getenv("RDS_PASSWORD"),
		os.Getenv("RDS_SSLMODE"),
		"ws://119.8.146.229:8546",
		os.Getenv("ORIGIN_ALLOWED"),
		"0",
	}
}
