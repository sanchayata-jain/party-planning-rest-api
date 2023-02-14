package tables

import (
	"context"
	"strconv"

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
	err := s.repository.CreateTable(ctx, table)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) CountNumberOfEmptySeats() (string, error) {
	emptySeats, err := s.repository.CountNumberOfEmptySeats()
	if err != nil {
		return "", nil
	}
	emptySeatsSum := strconv.Itoa(emptySeats)

	return emptySeatsSum, nil
}
