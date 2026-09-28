package handlers

import (
	"kolu/internal/repositories"
	"kolu/internal/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateNoteHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		var input types.CreateNoteRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		note, err := repositories.CreateNote(pool, &input, userId)

		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"data": &note})
	}

}

func GetAllNotesHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)

		notes, err := repositories.GetAllNotes(pool, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNoContent, gin.H{"error": "No Notes Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": notes})
	}
}

func GetAllTrashedNotesHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)

		notes, err := repositories.GetAllTrashedNotes(pool, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNoContent, gin.H{"error": "No Notes Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": notes})
	}
}

func GetNoteById(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid Note ID",
			})
			return
		}

		note, err := repositories.GetNoteByID(pool, id, userId)
		if err != nil {

			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Note Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": &note})

	}
}

func UpdateNote(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		id, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid Note ID",
			})
			return
		}

		var input types.UpdateNoteRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Title == nil && input.Content == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one field must be provided"})
			return
		}
		existing, err := repositories.GetNoteByID(pool, id, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Note Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		var title string
		if input.Title != nil {
			title = *input.Title
		} else {
			title = existing.Title
		}

		var content string
		if input.Content != nil {
			content = *input.Content
		} else {
			content = existing.Content
		}

		updateInput := types.UpdateNoteRequest{
			Title:   &title,
			Content: &content,
		}

		note, err := repositories.UpdateNote(pool, &updateInput, id, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Note Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": &note})
	}
}

func DeleteNoteById(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		id := ctx.Param("id")

		err := repositories.DeleteNoteById(pool, id, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Note Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func DeleteAllNotes(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)

		err := repositories.DeleteAllNotes(pool, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Note Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}
