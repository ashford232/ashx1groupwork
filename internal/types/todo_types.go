package types

import (
	"time"

	"github.com/google/uuid"
)

type CreateTodoRequest struct {
	ID          *uuid.UUID `json:"id"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	Priority    *string    `json:"priority"`
	DueAt       *time.Time `json:"due_at"`
}
type UpdateTodoRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	IsCompleted   *bool      `json:"is_completed"`
	IsDeleted   bool       `json:"is_deleted"`
	Priority    *string    `json:"priority"`
	DueAt       *time.Time `json:"due_at"`
}

