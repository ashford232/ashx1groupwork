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

func CreateDocHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		var input types.CreateDocRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		note, err := repositories.CreateDoc(pool, &input, userId)

		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"data": &note})
	}

}

func GetAllDocsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)

		docs, err := repositories.GetAllDocs(pool, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNoContent, gin.H{"error": "No Documents Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": docs})
	}
}

func GetAllTrashedDocsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)

		docs, err := repositories.GetAllTrashedDocs(pool, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNoContent, gin.H{"error": "No Documents Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": docs})
	}
}

func GetDocById(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid Document ID",
			})
			return
		}

		doc, err := repositories.GetDocByID(pool, id, userId)
		if err != nil {

			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Document Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": &doc})

	}
}

func UpdateDoc(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		id, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid Document ID",
			})
			return
		}

		var input types.UpdateDocRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		existing, err := repositories.GetDocByID(pool, id, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Document Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		var name string
		if input.Name != "" {
			name = input.Name
		} else {
			name = existing.Name
		}

		updateInput := types.UpdateDocRequest{
			Name: name,
		}

		doc, err := repositories.UpdateDoc(pool, &updateInput, id, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Document Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": &doc})
	}
}

func DeleteDocById(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)
		id := ctx.Param("id")

		err := repositories.DeleteDocById(pool, id, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Document Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func DeleteAllDocs(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		userId := GetUserID(ctx)

		err := repositories.DeleteAllDocs(pool, userId)
		if err != nil {
			log.Println(err)
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Document Not Found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}
