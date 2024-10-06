package middleware

import (
	"net/http"

	"github.com/ekefan/pre-sales-deal-tracker/backend/token"
	"github.com/gin-gonic/gin"
)

const (
	adminRole = "admin"
	salesRole = "sales"
)

func AdminAccessAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := ctx.MustGet(AuthPayloadKey).(*token.Payload)
		if payload.Role != adminRole {
			abortHandlingRequest(ctx, http.StatusForbidden, "UNAUTHORIZED",
				"user is not allowed to access this resource",
				"this resouces is restricted to admin users")
			return
		}
	}
}

func SalesAccessAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := ctx.MustGet(AuthPayloadKey).(*token.Payload)
		if payload.Role != salesRole {
			abortHandlingRequest(ctx, http.StatusForbidden, "UNAUTHORIZED",
				"user is not allowed to access this resource",
				"this resouces is restricted to sales users")
			return
		}
	}
}
