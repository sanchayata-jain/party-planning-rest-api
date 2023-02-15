package tables

import (
	"context"
	"encoding/json"

	// "net/http"

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
	createTable := &models.Table{}
	tx := r.db.Select("Capacity", "SeatsEmpty").Create(&table)
	if tx.Error != nil {
		return tx.Error
	}
	_, err := json.Marshal(*createTable)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) CountNumberOfEmptySeats() (int, error) {
	var count int
	tx := r.db.Model(&models.Table{}).Select("SUM(seats_empty)").Scan(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return count, nil
}

func EditEmptySeatsAfterGuestsLeave(capacity int, tableID int, db *gorm.DB) error {
	//using guests table number make empty seats equal capacity
	tx := db.Model(&models.Table{}).Select("seats_empty").Where("id = ?", tableID).Update("seats_empty", capacity)

	return tx.Error
}

func EditEmptySeatsAfterGuestsArrive(db *gorm.DB, emptySeats int, tableID int) error {
	tx := db.Model(&models.Table{}).Select("seats_empty").Where("id = ?", tableID).Update("seats_empty", emptySeats)

	return tx.Error
}

func GetTableCapacity(requestedTable int, db *gorm.DB) (int, error) {
	//if requestedTableId found, get capacity and return, otherwise return error (table not found)
	table := models.Table{}
	tx := db.First(&table, "ID = ?", requestedTable)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return table.Capacity, nil
}
