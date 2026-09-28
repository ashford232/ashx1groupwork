package handlers

import (
	"kolu/internal/repositories"
	"kolu/internal/types"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		var input types.CreateTodoRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repositories.CreateTodo(pool, &input, userId)

		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"data": &todo})
	}

}

func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)

		todos, err := repositories.GetAllTodos(pool, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNoContent, gin.H{"error": "No Todos Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": todos})
	}
}

func GetAllTrashedTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)

		todos, err := repositories.GetAllTrashedTodos(pool, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNoContent, gin.H{"error": "No Todos Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": todos})
	}
}

func GetTodoById(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid todo ID",
			})
			return
		}

		todo, err := repositories.GetTodoByID(pool, id, userId)
		if err != nil {

			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": &todo})

	}
}

func UpdateTodo(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		id, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid todo ID",
			})
			return
		}

		var input types.UpdateTodoRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Title == nil && input.Description == nil && input.IsCompleted == nil && input.DueAt == nil && input.Priority == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one field must be provided"})
			return
		}
		existing, err := repositories.GetTodoByID(pool, id, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		var isCompleted bool
		if input.IsCompleted != nil {
			isCompleted = *input.IsCompleted
		} else {
			isCompleted = existing.IsCompleted
		}

		var title string
		if input.Title != nil {
			title = *input.Title
		} else {
			title = existing.Title
		}

		var description string
		if input.Description != nil {
			description = *input.Description
		} else {
			description = existing.Description
		}

		var priority string
		if input.Priority != nil {
			priority = *input.Priority
		} else {
			priority = existing.Priority
		}

		var dueAt *time.Time
		if input.DueAt != nil {
			dueAt = input.DueAt
		} else {
			dueAt = existing.DueAt
		}

		updateInput := types.UpdateTodoRequest{
			Title:       &title,
			Description: &description,
			IsCompleted:   &isCompleted,
			Priority:    &priority,
			DueAt:       dueAt,
		}

		todo, err := repositories.UpdateTodo(pool, &updateInput, id, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": &todo})
	}
}

func DeleteTodoById(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		id := ctx.Param("id")

		err := repositories.DeleteTodoById(pool, id, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func DeleteAllTodos(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)

		err := repositories.DeleteAllTodos(pool, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}
