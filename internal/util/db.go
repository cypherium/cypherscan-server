package util

import (
	"fmt"

	"github.com/jet/go-interstellar"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// DbRunner is a interface mainly for a tests
type DbRunner interface {
	Run(RunFunc) error
}

// DbClient is a Database access client
type DbClient struct {
	db             *gorm.DB
	cosmosdbClient interstellar.Client
}

// Close close db connection
func (dbClient *DbClient) Close() {
	dbClient.db.Close()
}

// Run accept a function which take RunFunc as a parameter
func (dbClient *DbClient) Run(f RunFunc) error {
	return f(dbClient.db)
}

// ConnectDb will return a open db connection
func ConnectDb(drive string, args ...interface{}) (*DbClient, error) {
	if drive != "postgres" && drive != "sqlite3" {
		return nil, &MyError{Message: fmt.Sprintf("Unsupported db: %s, only supporting sqlite3 and postgres", drive)}
	}
	// name, port, dbName, userName, password := args
	connectionStr := ""
	if drive == "postgres" {
		connectionStr = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", args...)
		log.Info("connectionStr", fmt.Sprintf("%v", connectionStr))
	} else {
		connectionStr = args[2].(string)
		log.Info("connectionStr", fmt.Sprintf("%s", connectionStr))
	}
	_db, err := gorm.Open(drive, connectionStr)
	if err != nil {
		return nil, err
	}
	return &DbClient{db: _db}, nil
}

// RunFunc is the type of function to run db scripts
type RunFunc func(db *gorm.DB) error
