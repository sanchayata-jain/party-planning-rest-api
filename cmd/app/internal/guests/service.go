package guests

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/getground/tech-tasks/backend/cmd/app/internal/models"
	"github.com/getground/tech-tasks/backend/cmd/app/internal/tables"
)

type guestInfoToReturn struct {
	Name               string `json:"name"`
	Table              int    `json:"table"`
	AccompanyingGuests int    `json:"accompanying_guests"`
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) AddGuestToGuestList(ctx context.Context, guest models.Guest) error {
	// we are just adding the guest to the guest list, the guest has not actually arrived to the party yet
	// so we can leave TimeArrived empty
	guest.LeftParty = false
	// get capacity of table they wish to book
	capacity, err := tables.GetTableCapacity(guest.Table, s.repository.db)
	if err != nil {
		return err
	}

	if capacity < guest.AccompanyingGuests+1 {
		// table they have requested is not big enough so turn away the group
		return errors.New("the table you have requested does not have enough space for your group, try a different table")
	}
	//TODO: check if table is not already assinged using guest list

	// if guest.AccompanyingGuests + 1 >= guest.Table
	err = s.repository.AddGuestToGuestlist(ctx, guest)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetGuestsOnGuestList(ctx context.Context) ([]byte, error) {
	guests, err := s.repository.GetGuestsOnGuestList(ctx)
	guestInfo := []*models.Guest{}
	json.Unmarshal(guests, &guestInfo)
	if err != nil {
		return nil, err
	}

	guestList := []guestInfoToReturn{}

	for _, g := range guestInfo {
		guestList = append(guestList, guestInfoToReturn{
			Name:               g.Name,
			Table:              g.Table,
			AccompanyingGuests: g.AccompanyingGuests,
		})
	}

	b, err := json.Marshal(guestList)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s Service) EditGuestsList(ctx context.Context, guest models.Guest) error {
	//get name from path and then loop through the guest list to get tableID
	tableID, err := GetGuestTableID(ctx, s.repository.db, guest.Name)
	if err != nil {
		//    http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	// using tableID, loop through the tables, and see if the id matches any tables
	capacity, err := tables.GetTableCapacity(tableID, s.repository.db)
	if err != nil {
		return err
	}
	if capacity < guest.AccompanyingGuests + 1 {
		// table they have requested is not big enough so turn away the group
		return errors.New("the table you have requested does not have enough space for your new group size. goodbye")
	}
	// okay so capacity is fine .. lets update party size on guest list and add arrival time
	// csall repo
	arrivalTime := time.Now()
	err = s.repository.EditGuestList(arrivalTime, guest)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) DeleteGuestFromList(name string) error {
	err := s.repository.DeleteGuest(name)
	if err != nil {
		return err
	}

	return nil
} 