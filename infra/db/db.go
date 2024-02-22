package db

import (
	"io"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase opens DB connection and also migrate
func InitDatabase(dbName string) (*gorm.DB, error) {
	newLogger := logger.New(log.New(io.Discard, "\r\n", log.LstdFlags), logger.Config{})

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	// Get generic database object sql.DB to use its functions

	err = db.AutoMigrate(&Transaction{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
