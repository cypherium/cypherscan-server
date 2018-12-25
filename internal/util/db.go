package util

import (
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
	_db, err := gorm.Open(drive, args...)
	if err != nil {
		return nil, err
	}
	return &DbClient{db: _db}, nil
}

// RunFunc is the type of function to run db scripts
type RunFunc func(db *gorm.DB) error
