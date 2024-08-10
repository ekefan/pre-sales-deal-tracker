package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ekefan/deal-tracker/internal/token"
	"github.com/gin-gonic/gin"
)

const (
	authHeaderKey = "Authorization"
	authHeaderType = "bearer"
	authPayloadKey = "authorization_payload"
)

func authMiddleware(maker token.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Check header for authorization header
		authHeader := ctx.GetHeader(authHeaderKey); 
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("no auth header passed")))
			return
		}

		//check for format validation
		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid or missing bearer token")))
			return
		}

		//check authorization type
		authType := fields[0]
		if strings.ToLower(authType) != authHeaderType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("authorization type not supported")))
			return
		}
		
		//token verification
		payload, err := maker.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid access token: %v", err)))
			return
		}
		//if verification is successful pass to the next handler func
		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}