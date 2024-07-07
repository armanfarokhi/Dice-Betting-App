package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
