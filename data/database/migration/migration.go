package migration

import (
	"github.com/hashicorp/go-hclog"
	"github.com/vdbulcke/cert-manager/data"
	"gorm.io/gorm"
)

// DBMigration initialize DB
func DBMigration(db *gorm.DB, logger hclog.Logger) error {

	err := db.AutoMigrate(&data.Tag{}, &data.Certificate{})
	if err != nil {
		logger.Error("Error Running DB Migration", "error", err)
		return err
	}

	return nil
}
