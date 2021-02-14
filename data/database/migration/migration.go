package migration

import (
	"github.com/vdbulcke/cert-manager/data"
	"github.com/vdbulcke/cert-manager/database"
)

// DBMigration initialize DB
func DBMigration(db *database.DB) {

	// Create DB Model
	db.AutoMigrate(&data.Certificate{})
	db.AutoMigrate(&data.Tag{})
}
