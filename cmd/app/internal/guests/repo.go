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

func (r Repository) AddGuestToGuestlist(ctx context.Context, guest models.Guest) error {
	createGuest := &models.Guest{}
	tx := r.db.Select("Name", "Table", "AccompanyingGuests", "LeftParty").Create(&guest)
	if tx.Error != nil {
		return tx.Error
	}
	_, err := json.Marshal(*createGuest)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetGuestsOnGuestList(ctx context.Context) ([]byte, error) {
	guests := []*models.Guest{}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Guest{}).
			Find(&guests).
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

	b, err := json.Marshal(guests)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GetGuestTableID(ctx context.Context, db *gorm.DB, name string) (int, error) {
	guest := &models.Guest{}

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Guest{}).
			Find(&guest).
			Where("name = ?", name).
			WithContext(ctx).
			Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return guest.Table, nil
}

func (r Repository) EditGuestList(arrivalTime time.Time, guest models.Guest) error {
	tx := r.db.Model(&models.Guest{}).Select("time_arrived", "accompanying_guests").Where("name = ?", guest.Name).Updates(models.Guest{TimeArrived: arrivalTime, AccompanyingGuests: guest.AccompanyingGuests})
	return tx.Error
}

func (r Repository) DeleteGuest(name string) error {
	tx := r.db.Where("name = ?", name).Delete(&models.Guest{})

	return tx.Error
}

func (r Repository) GetArrivedGuests(ctx context.Context) ([]byte, error){
	arrivedGuests := []*models.Guest{}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Guest{}).
			Find(&arrivedGuests).
			Where("time_arrived <> ?", nil).
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

	b, err := json.Marshal(arrivedGuests)
	if err != nil {
		return nil, err
	}

	return b, nil
}
