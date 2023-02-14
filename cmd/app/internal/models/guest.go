package models

import "time"

type Guest struct {
	Name               string    `gorm:"name;PRIMARY_KEY"`
	Table              int       `gorm:"table"`
	AccompanyingGuests int       `gorm:"accompanying_guests"`
	TimeArrived        time.Time `gorm:"time_arrived"`
	LeftParty          bool      `gorm:"left_party"`
}
