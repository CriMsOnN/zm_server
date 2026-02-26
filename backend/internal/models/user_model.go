package models

import (
	"time"
)

type User struct {
	ID         string    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Identifier string    `json:"identifier" db:"identifier"`
	LastIP     string    `json:"last_ip" db:"last_ip"`
	LastLogin  time.Time `json:"last_login" db:"last_login"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
