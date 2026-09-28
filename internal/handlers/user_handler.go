package handlers

import (
	"errors"
	"kolu/internal/config"
	errork "kolu/internal/error-k"
	"kolu/internal/repositories"
	"kolu/internal/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SignUpUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.SignUpUserRequest
		err := ctx.ShouldBindJSON(&input)
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		user, err := repositories.SignUp(pool, &input)
		if err != nil {
			log.Println(err)
			if isUniqueViolation(err) {
				ctx.JSON(http.StatusConflict, gin.H{"error": errork.ErrUserExists.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, errork.ErrInternalServerError)
			return

		}
		ctx.JSON(http.StatusCreated, gin.H{"user": &user})

	}
}
func SignInUserHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.SingInUserRequest
		err := ctx.ShouldBindJSON(&input)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errork.ErrInvalidRequestBody)
			return
		}
		userAndToken, err := repositories.SignIn(pool, &input, cfg)
		if err != nil {
			log.Println(err)
			if errors.Is(err, errork.ErrInvalidCredentials) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": errork.ErrInvalidCredentials.Error()})
				return
			}

			ctx.JSON(http.StatusInternalServerError, errork.ErrInternalServerError)
			return
		}
		
		ctx.JSON(http.StatusOK,&userAndToken)
	}
}

func GetUserByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		user, err := repositories.GetUserByID(pool, id)
		if err != nil {
			if errors.Is(err, errork.ErrInvalidCredentials) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": errork.ErrInvalidCredentials})

				return
			}

			ctx.JSON(http.StatusInternalServerError, errork.ErrInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"user": &user})

	}
}



func GetUser(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id := GetUserID(ctx)

		user, err := repositories.GetUserByID(pool, id)
		if err != nil {
			if errors.Is(err, errork.ErrInvalidCredentials) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": errork.ErrInvalidCredentials})

				return
			}

			ctx.JSON(http.StatusInternalServerError, errork.ErrInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"user": &user})

	}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
