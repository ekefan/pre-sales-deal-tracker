package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ekefan/pre-sales-deal-tracker/backend/middleware"
	"github.com/ekefan/pre-sales-deal-tracker/backend/token"
	"github.com/gin-gonic/gin"
)

type CreateUserReq struct {
	Username string `json:"username" binding:"required,gte=4,lte=6,alphanum"`
	Fullname string `json:"fullname" binding:"required,gte=4"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=admin sales manager"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserReq
	if err := bindClientRequest(ctx, &req, jsonSource); err != nil {
		slog.Error(err.Error())
		return
	}

	payload := ctx.MustGet(middleware.AuthPayloadKey)
	payloadData := payload.(*token.Payload)
	fmt.Printf("%T, %v, %v\n", payloadData, payloadData.Role, req.Role)
	ctx.JSON(http.StatusCreated, successMessage("success"))
}


