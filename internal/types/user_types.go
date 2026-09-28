package types

import "kolu/internal/models"

type SignUpUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" `
	Email    string `json:"email" binding:"required" `
	Password string `json:"password"  binding:"required"`
}

type SingInUserRequest struct {
	Email    string `json:"email" binding:"required" `
	Password string `json:"password"  binding:"required"`
}

type SignInUserResponse struct{
	Token string `json:"token"`
	User *models.User `json:"user"`
}
type SetupUserRequest struct{
	Username string `json:"username" binding:"required"`
	PhotoURL string `json:"photo_url"`
}
type PresignUploadRequest struct {
	Filename string `json:"filename" binding:"required"`
	MimeType string `json:"mime_type" binding:"required"`
}

type PresignUploadResponse struct {
	UploadURL string `json:"upload_url"`
	PublicURL string `json:"public_url"`
}