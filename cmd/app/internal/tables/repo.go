package tables

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"github.com/getground/tech-tasks/backend/cmd/app/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) GetTables(ctx context.Context) ([]byte, error) {
	tables := []*models.Table{}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Table{}).
			Find(&tables).
			WithContext(ctx).
			Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(tables)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r Repository) CreateTable(ctx context.Context, table models.Table) error {
	err := r.db.Select("Capacity", "SeatsEmpty").Create(&table).Error
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetLastTableMade() ([]byte, error) {
	createTable := &models.Table{}
	err := r.db.Last(&createTable).Error
	if err != nil {
		return nil, err
	}
	type tableInfoToReturn struct {
		ID       int `json:"id"`
		Capacity int `json:"capacity"`
	}
	info := tableInfoToReturn {
		ID: createTable.ID,
		Capacity: createTable.Capacity,
	}

	b, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r Repository) CountNumberOfEmptySeats() (int, error) {
	var count int
	tx := r.db.Debug().Model(&models.Table{}).Select("SUM(seats_empty)").Scan(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return count, nil
}

func EditEmptySeatsAfterGuestsLeave(capacity int, tableID int, db *gorm.DB) error {
	tx := db.Model(&models.Table{}).Select("seats_empty").Where("id = ?", tableID).Update("seats_empty", capacity)

	return tx.Error
}

func EditEmptySeatsAfterGuestsArrive(db *gorm.DB, emptySeats int, tableID int) error {
	err := db.Model(&models.Table{}).Select("seats_empty").Where("id = ?", tableID).Update("seats_empty", emptySeats).Error
	if err != nil {
		return err
	}
	return nil
}

func GetTableCapacity(requestedTable int, db *gorm.DB) (int, error) {
	table := models.Table{}
	tx := db.Debug().First(&table, "ID = ?", requestedTable)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return table.Capacity, nil
}
