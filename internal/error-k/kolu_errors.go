package errork

import (
	"github.com/gin-gonic/gin"
)

var (
	ErrInternalServerError = gin.H{"error": "Internal Server Error"}
	ErrNoAuthorizationError     = gin.H{"error": "Authorization required"}
	ErrInvalidAuthorizationError     = gin.H{"error": "Invalid Authorization"}
	ErrInvalidTokenError     = gin.H{"error": "Invalid Token"}
	ErrInvalidTokenClaimsError     = gin.H{"error": "Token Expired"}
	ErrTokenExpiredError    = gin.H{"error": "Invalid Token Claims"}
	ErrInvalidRequestBody  = gin.H{"error": "Invalid Request Body"}
	ErrInvalidParamBody    = gin.H{"error": "Invalid Param Body"}
)
