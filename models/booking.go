package models

import "time"

type Booking struct {
	ID            uint   `gorm:"primaryKey"`
	UserID        uint   `json:"user_id"`
	PassangerName string `json:"passenger_name"`
	FlightID      string `json:"flight_id"`
	FlightName    string `json:"flight_name"`
	Departure     string `json:"flight_departure"`
	Arrival       string `json:"flight_arrival"`
	SeatNumber    string `json:"seat_number"`
	Date          string `json:"date"`
	CreatedAt     time.Time
}
