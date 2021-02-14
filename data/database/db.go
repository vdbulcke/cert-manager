package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB generic db
type DB struct {
	gorm.DB
}

// NewSqliteDB create a sqlite db connection
func NewSqliteDB(dbpath string) *gorm.DB {

	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db

}
