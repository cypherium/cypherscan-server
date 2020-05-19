package bizutil

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config hold the all Program's configurations comming from .env file
type Config struct {
	ExecutionTimeout  int
	NodesUrls         []string
	DynamoDbRegion    string
	RecentTTLDuration time.Duration
	OriginAllowed     string
}

// GetConfig will get the config values from .env file
func GetConfig() (*Config, error) {
	godotenv.Load()
	executionTimeout, err := strconv.Atoi(os.Getenv("EXECUTION_TIMEOUT"))
	if err != nil {
		return nil, err
	}
	recentTTLDuration, err := strconv.Atoi(os.Getenv("RECENT_TTL_DURATION_IN_SECONDS"))
	if err != nil {
		return nil, err
	}

	return &Config{
		ExecutionTimeout:  executionTimeout,
		NodesUrls:         strings.Split(os.Getenv("NODES_URLS"), ","),
		DynamoDbRegion:    os.Getenv("DYNAMODB_REGION"),
		RecentTTLDuration: time.Duration(recentTTLDuration * 1e9),
		OriginAllowed:     os.Getenv("ORIGIN_ALLOWED"),
	}, nil
}
