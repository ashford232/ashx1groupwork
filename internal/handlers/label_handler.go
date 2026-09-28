package handlers

import (
	"errors"
	errork "kolu/internal/error-k"
	"kolu/internal/repositories"
	"kolu/internal/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllUserLabels(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		user_id := GetUserID(ctx)

		labels, err := repositories.GetUsersLabels(pool, user_id)
		if err != nil {
			if errors.Is(err, errork.ErrInvalidCredentials) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": errork.ErrInvalidCredentials.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, errork.ErrInternalServerError)

			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": &labels})

	}
}

func CreateLabelHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)

		var input types.CreateLabelRequest

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		label, err := repositories.CreateLabel(pool, &input, userId)
		if err != nil {
			log.Println(err)

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"data": &label})
	}
}

func GetLabelHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		labelID := ctx.Param("id")

		label, err := repositories.GetLabel(pool, labelID, userId)
		if err != nil {
			log.Println(err)

			if errors.Is(err, pgx.ErrNoRows) {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "Label Not Found",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": &label})
	}
}

func UpdateLabelHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		labelID := ctx.Param("id")

		var input types.UpdateLabelRequest

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		label, err := repositories.UpdateLabel(
			pool,
			labelID,
			&input,
			userId,
		)

		if err != nil {
			log.Println(err)

			if errors.Is(err, pgx.ErrNoRows) {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "Label Not Found",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": &label})
	}
}

func DeleteLabelHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		labelID := ctx.Param("id")

		err := repositories.DeleteLabel(
			pool,
			labelID,
			userId,
		)

		if err != nil {
			log.Println(err)

			if errors.Is(err, pgx.ErrNoRows) {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "Label Not Found",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}
func AddTodoLabelHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)

		var input types.AddTodoLabelRequest

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		_, err := uuid.Parse(input.LabelID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Label ID",
			})

			return
		}
		_, err = uuid.Parse(input.TodoID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Todo ID",
			})
			return
		}
		label, err := repositories.AddTodoLabels(pool, &input, userId)
		if err != nil {
			log.Println(err)
			if isUniqueViolation(err) {
				ctx.JSON(http.StatusConflict, gin.H{
					"error": "Label already exists",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"data": label})
	}
}

func AddNoteLabelHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)

		var input types.AddNoteLabelRequest

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		_, err := uuid.Parse(input.LabelID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Note ID",
			})

			return
		}
		_, err = uuid.Parse(input.NoteID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Note ID",
			})
			return
		}
		label, err := repositories.AddNoteLabels(pool, &input, userId)
		if err != nil {
			log.Println(err)
			if isUniqueViolation(err) {
				ctx.JSON(http.StatusConflict, gin.H{
					"error": "Label already exists",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"data": label})
	}
}
func AddDocLabelHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)

		var input types.AddDocLabelRequest

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		_, err := uuid.Parse(input.LabelID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Document ID",
			})

			return
		}
		_, err = uuid.Parse(input.DocID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Document ID",
			})
			return
		}
		label, err := repositories.AddDocLabels(pool, &input, userId)
		if err != nil {
			log.Println(err)
			if isUniqueViolation(err) {
				ctx.JSON(http.StatusConflict, gin.H{
					"error": "Label already exists",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"data": label})
	}
}
