package models

import "time"

type Guest struct {
	Name               string    `gorm:"PRIMARY_KEY"`
	Table              int       `json:"table"`
	AccompanyingGuests int       `json:"accompanying_guests"`
	TimeArrived        time.Time `json:"time_arrived"`
	LeftParty          bool      `json:"left_party"`
}
