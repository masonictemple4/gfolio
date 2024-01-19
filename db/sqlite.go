package db

import (
	"fmt"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Will Panic if it can't connect to the database
// If name does not contain the .db extension we will add it.
func NewSqliteDB(name string, config *gorm.Config) *gorm.DB {
	if filepath.Ext(name) != ".db" {
		name = fmt.Sprintf("%s.db", name)
	}

	db, err := gorm.Open(sqlite.Open(name))
	if err != nil {
		panic(fmt.Errorf("newsqlitedb: failed to connect to database: %w", err))
	}

	return db
}
