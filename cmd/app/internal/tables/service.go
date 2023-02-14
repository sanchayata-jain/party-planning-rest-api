package tables

import (
	"context"

	"github.com/getground/tech-tasks/backend/cmd/app/internal/models"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) GetTables(ctx context.Context) ([]byte, error) {
	tables, err := s.repository.GetTables(ctx)
	if err != nil {
		return nil, err
	}

	return tables, nil
}

func (s Service) CreateTable(ctx context.Context, table models.Table) error {
	table.SeatsEmpty = table.Capacity
	// table.ID = 1
	err := s.repository.CreateTable(ctx, table)
	if err != nil {
		return err
	}
	return nil
}
