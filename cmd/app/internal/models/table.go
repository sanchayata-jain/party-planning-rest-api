package models

// Table is a struct which holds the information about a table
type Table struct {
	ID         int `gorm:"id;AUTO_INCREMENT"`
	Capacity   int `gorm:"capacity"`
	SeatsEmpty int `gorm:"seats_empty"`
}
