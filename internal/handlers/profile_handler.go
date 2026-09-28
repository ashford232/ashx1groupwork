package handlers

import (
	"log"
	"net/http"

	"kolu/internal/config"
	errork "kolu/internal/error-k"
	"kolu/internal/repositories"
	"kolu/internal/storage"
	"kolu/internal/types"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPresignedAvatarUploadHandler(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.PresignUploadRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, errork.ErrInvalidRequestBody)
			return
		}

		userID := GetUserID(ctx) // Extracted from JWT context via middleware

		r2Storage, err := storage.NewR2Storage(ctx, &cfg.Cloudflare)
		if err != nil {
			log.Println("Failed to initialize R2 storage:", err)
			ctx.JSON(http.StatusInternalServerError, errork.ErrInternalServerError)
			return
		}

		uploadURL, publicURL, err := r2Storage.GeneratePresignedUpload(ctx, userID, input.Filename, input.MimeType)
		if err != nil {
			log.Println("Failed to generate presigned URL:", err)
			ctx.JSON(http.StatusInternalServerError, errork.ErrInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, types.PresignUploadResponse{
			UploadURL: uploadURL,
			PublicURL: publicURL,
		})
	}
}

func SetupProfileHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.SetupUserRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, errork.ErrInvalidRequestBody)
			return
		}

		userIDStr := GetUserID(ctx)
		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user session"})
			return
		}

		user, err := repositories.SetupProfile(pool, &input, userUUID)
		if err != nil {
			log.Println("SetupProfile error:", err)
			if isUniqueViolation(err) {
				ctx.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, errork.ErrInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"user": &user})
	}
}
