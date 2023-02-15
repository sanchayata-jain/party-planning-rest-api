package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	goModels "github.com/getground/tech-tasks/backend/cmd/app/internal/models"
)

func ReadyStateDB(databaseURL string, models ...interface{}) (*gorm.DB, error) {
	db, err := NewPostgresDatabase(databaseURL) // error here
	if err != nil {
		return nil, fmt.Errorf("error getting postgres connection: %w", err)
	}

	err = MigrateDbSchema(db, models...)
	if err != nil {
		return nil, fmt.Errorf("error migrating database schema: %w", err)
	}

	return db, nil
}

func MigrateDbSchema(db *gorm.DB, models ...interface{}) error {
	m := []interface{}{
		&goModels.Table{},
		&goModels.Guest{},
	}
	m = append(m, models...)
	return db.AutoMigrate(
		m...,
	)
}

func NewPostgresDatabase(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		DSN: databaseURL,
	}), &gorm.Config{})
}
