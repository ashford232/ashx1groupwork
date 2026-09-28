package models

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	IsCompleted bool       `json:"is_completed" db:"is_completed"`
	Priority    string     `json:"priority" db:"priority"`
	DueAt       *time.Time `json:"due_at" db:"due_at"`
	UserId      string     `json:"user_id" db:"user_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	IsDeleted   bool       `json:"is_deleted" db:"is_deleted"`
	Labels      []Label    `json:"labels" db:"-"`
}
