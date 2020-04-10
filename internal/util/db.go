package util

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

// DbRunner is a interface mainly for a tests
type DbRunner interface {
	Run(RunFunc) error
}

// DbClient is a Database access client
type DbClient struct {
	db *gorm.DB
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
		return nil, &MyError{fmt.Sprintf("Unsupported db: %s, only supporting sqlite3 and postgres", drive)}
	}
	// name, port, dbName, userName, password := args
	connectionStr := ""
	if drive == "postgres" {
		connectionStr = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", args...)
	} else {
		connectionStr = args[2].(string)
		fmt.Println("connectionStr", connectionStr)
	}
	_db, err := gorm.Open(drive, connectionStr)
	if err != nil {
		return nil, err
	}
	return &DbClient{db: _db}, nil
}

// RunFunc is the type of function to run db scripts
type RunFunc func(db *gorm.DB) error
