package models

import (
	"time"
	"github.com/google/uuid"

)

type Note struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Content string     `json:"content" db:"content"`
	UserId      string     `json:"user_id" db:"user_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted   bool       `json:"is_deleted" db:"is_deleted"`
	Labels      []Label    `json:"labels" db:"-"`
}

