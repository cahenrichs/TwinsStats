package db

import (
	"github.com/cahenrichs/mlbstats/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Migrate the schema
	err = db.AutoMigrate(
		&models.Player{},
		&models.Team{},
		&models.HittingStats{},
		&models.PitchingStats{},
	)
	return db, err
}
