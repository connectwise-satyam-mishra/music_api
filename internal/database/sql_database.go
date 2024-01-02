package database

import (
	"log"
	"music_api/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB connects to the database and performs auto migration for the models.
// It returns the database instance and any error encountered.
func ConnectDB() (*gorm.DB, error) {
	// Connect to the database
	db, err := gorm.Open(sqlite.Open("music_tracks.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.MusicTrack{})
	err = db.AutoMigrate(&models.Artist{})
	if err != nil {
		log.Fatal("Failed to migrate the database:", err)
	}
	DB = db // Store the database instance globally for use in handlers or models
	return db, nil
}

// CloseDB closes the database connection.
// It returns any error encountered while closing the connection.
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
