package guests

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/getground/tech-tasks/backend/cmd/app/internal/models"
	"github.com/getground/tech-tasks/backend/cmd/app/internal/tables"
)

type arrivedGuestInfoToReturn struct {
	Name               string    `json:"name"`
	AccompanyingGuests int       `json:"accompanying_guest"`
	TimeArrived        time.Time `json:"time_arrived"`
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) AddGuestToGuestList(ctx context.Context, guest models.Guest) (string, error) {
	capacity, err := tables.GetTableCapacity(guest.Table, s.repository.db)
	if err != nil {
		return "", err
	}

	if capacity < guest.AccompanyingGuests+1 {
		// table they have requested is not big enough so turn away the group
		return "", errors.New("the table you have requested does not have enough space for your group, try a different table")
	}
	//TODO: check if table is not already assinged using guest list
	guest.Arrived = false
	name, err := s.repository.AddGuestToGuestlist(ctx, guest)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (s Service) GetGuestsOnGuestList(ctx context.Context) ([]byte, error) {
	guests, err := s.repository.GetGuestsOnGuestList(ctx)
	guestInfo := []*models.Guest{}
	json.Unmarshal(guests, &guestInfo)
	if err != nil {
		return nil, err
	}

	var guestList []string

	for _, g := range guestInfo {
		guestList = append(guestList, g.Name)
	}

	b, err := json.Marshal(guestList)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s Service) EditGuestsList(ctx context.Context, guest models.Guest) error {
	tableID, err := GetGuestTableID(ctx, s.repository.db, guest.Name)
	if err != nil {
		return err
	}
	capacity, err := tables.GetTableCapacity(tableID, s.repository.db)
	if err != nil {
		return err
	}
	if capacity < guest.AccompanyingGuests+1 {
		// table they have requested is not big enough so turn away the group
		return errors.New("the table you have requested does not have enough space for your new group size. goodbye")
	}
	guest.Arrived = true
	arrivalTime := time.Now()
	err = s.repository.EditGuestList(arrivalTime, guest)
	if err != nil {
		return err
	}
	updatedEmptySeats := capacity - (guest.AccompanyingGuests + 1)
	err = tables.EditEmptySeatsAfterGuestsArrive(s.repository.db, updatedEmptySeats, tableID)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) DeleteGuestFromList(ctx context.Context, name string) error {
	guest, err := s.repository.GetGuestFromName(name)
	if err != nil {
		return err
	}
	if !guest.Arrived {
		return errors.New("cannot checkout a guest that has not checked in yet")
	}
	
	tableID, err := GetGuestTableID(ctx, s.repository.db, name)
	if err != nil {
		return err
	}
	
	capacity, err := tables.GetTableCapacity(tableID, s.repository.db)
	if err != nil {
		return err
	}
	err = tables.EditEmptySeatsAfterGuestsLeave(capacity, tableID, s.repository.db)
	if err != nil {
		return err
	}
	err = s.repository.DeleteGuest(name)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetArrivedGuests(ctx context.Context) ([]byte, error) {
	arrivedGuests, err := s.repository.GetArrivedGuests(ctx)
	if err != nil {
		return nil, err
	}

	guestInfo := []*models.Guest{}
	err = json.Unmarshal(arrivedGuests, &guestInfo)
	if err != nil {
		return nil, err
	}

	arrivedGuestList := []arrivedGuestInfoToReturn{}

	for _, g := range guestInfo {
		arrivedGuestList = append(arrivedGuestList, arrivedGuestInfoToReturn{
			Name:               g.Name,
			AccompanyingGuests: g.AccompanyingGuests,
			TimeArrived:        g.TimeArrived,
		})
	}

	b, err := json.Marshal(arrivedGuestList)
	if err != nil {
		return nil, err
	}

	return b, nil
}
