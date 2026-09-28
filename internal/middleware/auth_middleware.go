package middleware

import (
	"fmt"
	"kolu/internal/config"
	errork "kolu/internal/error-k"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, errork.ErrNoAuthorizationError)
			ctx.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenStr == "" || tokenStr == authHeader {
			ctx.JSON(http.StatusUnauthorized, errork.ErrInvalidAuthorizationError)
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Method.Alg())
			}

			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errork.ErrInvalidTokenError)
			ctx.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, errork.ErrInvalidTokenClaimsError)
			ctx.Abort()
			return
		}

		userID, userIdOK := claims["user_id"].(string)
		email, emailOK := claims["email"].(string)
		exp, expOK := claims["exp"].(float64)

		if !userIdOK || !emailOK || !expOK {
			ctx.JSON(http.StatusUnauthorized, errork.ErrInvalidTokenClaimsError)
			ctx.Abort()
			return
		}

		expirationTime := time.Unix(int64(exp), 0)
		if time.Now().After(expirationTime) {
			ctx.JSON(http.StatusUnauthorized, errork.ErrTokenExpiredError)
			ctx.Abort()
			return
		}
		ctx.Set("user_id", userID)
		ctx.Set("email", email)
		ctx.Next()
	}
}
