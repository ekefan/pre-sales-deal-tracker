package middleware

import (
	"github.com/gin-gonic/gin"
)

// abortHandlingRequest takes the ctx prevents handlers from being called
// then sends json response to client
func abortHandlingRequest(ctx *gin.Context, httpStatus int, errCode, errMsg, details string) {
	ctx.AbortWithStatusJSON(httpStatus, gin.H{
		"code":  errCode,
		"error": errMsg,
		"details": details,
	})
}