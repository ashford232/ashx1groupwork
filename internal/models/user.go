package models

import "time"

type User struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Username  string `json:"username" db:"username"`
	Email     string `json:"email" db:"email"`
	PhotoURL     *string `json:"photo_url" db:"photo_url"`
	Password     string `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
