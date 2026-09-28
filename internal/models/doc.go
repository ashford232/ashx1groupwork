package models

import (
	"time"

	"github.com/google/uuid"
)

type Doc struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	ObjectKey string    `json:"object_key" db:"object_key"`
	MimeType  string    `json:"mime_type" db:"mime_type"`
	Size      int       `json:"size" db:"size"`
	UserId    string    `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
	Labels    []Label   `json:"labels" db:"-"`
}
