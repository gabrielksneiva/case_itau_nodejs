package connection

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewSqliteConnection returns a gorm DB instance
func NewSqliteConnection(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
