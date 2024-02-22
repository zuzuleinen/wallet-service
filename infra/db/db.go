package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDatabase opens DB connection and also migrate
func InitDatabase(dbName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(1000)

	err = db.AutoMigrate(&Transaction{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
