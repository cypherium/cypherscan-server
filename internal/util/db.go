package util

import (
	"github.com/jinzhu/gorm"

	"gitlab.com/ron-liu/cypherscan-server/internal/env"
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

// Connect will return a open db connection
func Connect() (*DbClient, error) {
	_db, err := gorm.Open(env.Env.DbDrive, env.Env.DbSource)
	if err != nil {
		return nil, err
	}
	return &DbClient{db: _db}, nil
}

// RunFunc is the type of function to run db scripts
type RunFunc func(db *gorm.DB) error
