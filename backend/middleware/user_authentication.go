package middleware

import (
	"net/http"
	"strings"

	"github.com/ekefan/pre-sales-deal-tracker/backend/token"
	"github.com/gin-gonic/gin"
)

const (
	authoHeaderKey       = "authorization"
	authHeaderTypeBearer = "bearer"
	AuthPayloadKey       = "auth_payload"
)

// UserAuthentication checks client requests to verify appropriate
// authorization header is used and access token is still valid
func UserAuthentication(tokenGenerator token.TokenGenerator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		if len(header) == 0 {
			abortHandlingRequest(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "authorization error", "no authorization header passed")
			return
		}

		headerFields := strings.Fields(header)
		if len(headerFields) != 2 {
			abortHandlingRequest(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "authorization error", "invalid authorization header")
			return
		}

		bearerField := strings.ToLower(headerFields[0])
		if bearerField != authHeaderTypeBearer {
			abortHandlingRequest(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "authorization error", "unsupported authorization type")
			return
		}

		accessToken := headerFields[1]
		payload, err := tokenGenerator.VerifyToken(accessToken)
		if err != nil {
			abortHandlingRequest(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "authorization error", "invalid access token")
			return
		}

		ctx.Set(AuthPayloadKey, payload)
		ctx.Next()
	}
}
