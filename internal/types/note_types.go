package types

import (
	"github.com/google/uuid"
)

type CreateNoteRequest struct {
	ID      *uuid.UUID `json:"id"`
	Title   string     `json:"title" binding:"required"`
	Content string     `json:"content"`
}
type UpdateNoteRequest struct {
	Title       *string `json:"title"`
	Content *string `json:"content"`
	IsDeleted   bool       `json:"is_deleted"`
}

