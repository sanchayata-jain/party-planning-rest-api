package server

import (
	"context"

	"github.com/getground/tech-tasks/backend/cmd/app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(ctx context.Context) (*gorm.DB, error) {
	psqlconn := "postgresql://user:password@localhost:5432/database?sslmode=disable"
	return gorm.Open(postgres.New(postgres.Config{
		DSN: psqlconn,
	}), &gorm.Config{})
}

func Init(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Table{},
		&models.Guest{},
	)
}
