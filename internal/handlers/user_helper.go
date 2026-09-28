package handlers

import "github.com/gin-gonic/gin"

func GetUserID(ctx *gin.Context) string {
	return ctx.MustGet("user_id").(string)
}