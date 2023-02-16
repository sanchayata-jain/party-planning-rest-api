package guests

import (
	"context"
	"encoding/json"
	"time"

	"github.com/getground/tech-tasks/backend/cmd/app/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) AddGuestToGuestlist(ctx context.Context, guest models.Guest) (string, error) {
	err := r.db.Debug().Select("Name", "Table", "AccompanyingGuests").Create(&guest).Error
	if err != nil {
		return "", err
	}

	return guest.Name, nil
}

func (r Repository) GetGuestsOnGuestList(ctx context.Context) ([]byte, error) {
	guests := []*models.Guest{}

	err := r.db.Model(&models.Guest{}).
		Find(&guests).
		WithContext(ctx).
		Error
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(guests)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GetGuestTableID(ctx context.Context, db *gorm.DB, name string) (int, error) {
	guest := &models.Guest{}
	err := db.Model(&models.Guest{}).
		Find(&guest).
		Where("name = ?", name).
		// WithContext(ctx).
		Error
	if err != nil {
		return 0, err
	}

	return guest.Table, nil
}

func (r Repository) EditGuestList(arrivalTime time.Time, guest models.Guest) error {
	tx := r.db.Debug().Model(&models.Guest{}).Select("time_arrived", "accompanying_guests").Where("name = ?", guest.Name).Updates(models.Guest{TimeArrived: arrivalTime, AccompanyingGuests: guest.AccompanyingGuests})
	return tx.Error
}

func (r Repository) DeleteGuest(name string) error {
	tx := r.db.Where("name = ?", name).Delete(&models.Guest{})

	return tx.Error
}

func (r Repository) GetArrivedGuests(ctx context.Context) ([]byte, error){
	arrivedGuests := []*models.Guest{}
	err := r.db.Debug().Model(&models.Guest{}).
		Find(&arrivedGuests).
		Where("time_arrived <> NULL").
		WithContext(ctx).
		Error

	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(arrivedGuests)
	if err != nil {
		return nil, err
	}

	return b, nil
}
