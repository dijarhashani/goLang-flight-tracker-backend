package models

import "time"

type Favorites struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `json:"user_id"`
	AirportIATA string `json:"airport_iata"`
	AirportName string `json:"airport_name"`
	CreatedAt   time.Time
}
