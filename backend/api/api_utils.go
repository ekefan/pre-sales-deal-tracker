package api

import "github.com/gin-gonic/gin"

func errorResponse(err error, code string) gin.H {
	return gin.H{
		"code": code,
		"error": err.Error(),
	}
}

func successMessage(msg string) gin.H {
	return gin.H{
		"message": msg,
	}
}