package infrastructure

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

	err = db.AutoMigrate(&Transaction{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
