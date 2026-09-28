package types

import (
	"github.com/google/uuid"
)

type CreateDocRequest struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ObjectKey string    `json:"object_key"`
	MimeType  string    `json:"mime_type"`
	Size      int       `json:"size"`
}
type UpdateDocRequest struct {
	Name      string    `json:"name"`
	IsDeleted   bool       `json:"is_deleted"`
}

