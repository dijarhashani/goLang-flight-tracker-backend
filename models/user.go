package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email     string `gorm:"unique" json:"email"`
	Password  string `json:"-"`
	CreatedAt time.Time
}
