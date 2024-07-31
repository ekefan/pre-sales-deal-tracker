package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateUsrReq holds fields that must be provided by client to create user
type CreateUsrReq struct {
	Username  string    `json:"username" binding:"required,alphanum"`
	Role      string    `json:"Role" binding:"required"`
	FullName  string    `json:"fullname" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
}

type CreateUsrResp struct {
	Username  string    `json:"username"`
	Role      string    `json:"Role"`
	FullName  string    `json:"fullname"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// adminCreateUserHandler http handler for the api end point for creating a new user
func (s *Server) adminCreateUserHandler(ctx *gin.Context) {
	var req CreateUsrReq
	//validate and bind request
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//create args to call data from the database the make call and return response depending
}

// adminUpdateUserHandler http handler for the api end point for updating a user
func (s *Server) adminUpdateUserHandler(ctx *gin.Context) {}

// adminDeleteUserhandler http handler for the api end point for Deleting a user
func (s *Server) adminDeleteUserHandler(ctx *gin.Context) {}

// adminCreateDealHandler http handler for the api end point for creating a new deal
func (s *Server) adminCreateDealHandler(ctx *gin.Context) {}

// adminUpdateDealsHandler http handler for the api end point for updating a deal
func (s *Server) adminUpdateDealHandler(ctx *gin.Context) {}

// adminDeleteDealHandler http handler for the api end point for Deleting a deal
func (s *Server) adminDeleteDealHandler(ctx *gin.Context) {}

// listUsershandler http handler for the api end point for getting list of users currently
func (s *Server) listUsersHandler(ctx *gin.Context) {}



// ================TODO=============================
//1. Create Custom validation for role: oneof - admin, sales, manager