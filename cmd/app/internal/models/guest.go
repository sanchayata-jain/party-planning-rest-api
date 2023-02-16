package models

import "time"

type Guest struct {
	Name               string    `gorm:"PRIMARY_KEY"`
	Table              int       `gorm:"foreignKey:Table"`
	AccompanyingGuests int       `json:"accompanying_guests"`
	TimeArrived        time.Time `json:"time_arrived"`
	Arrived            bool      `json:"arrived"`
}
