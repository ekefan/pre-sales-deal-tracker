package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	db "github.com/ekefan/pre-sales-deal-tracker/backend/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// CreateUserReq holds fields needed to create a user resource
type CreateUserReq struct {
	Username string `json:"username" binding:"required,gte=4,lte=6,alphanum"`
	Fullname string `json:"fullname" binding:"required,gte=4"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=admin sales manager"`
}

// createUser api handler for creating user resource
func (server *Server) createUsers(ctx *gin.Context) {
	var req CreateUserReq
	if err := bindClientRequest(ctx, &req, jsonSource); err != nil {
		slog.Error(err.Error())
		return
	}
	if !authAccess(ctx, []string{adminRole}) {
		return
	}
	hash, err := HashPassword(db.DefaultUserPassword)
	if err != nil {
		handleServerError(ctx, err)
		return
	}
	userID, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		FullName: req.Fullname,
		Email:    req.Email,
		Role:     req.Role,
		Password: hash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			slog.Warn("user already exist", "username", req.Username, "email", req.Email)
			statusCode = http.StatusConflict
			errCode = "STATUS_CONFLICT"
			errMsg = "can not create user"
			details = fmt.Sprintf("user exists with %v or %v", req.Email, req.Username)
			ctx.JSON(statusCode, errorResponse(statusCode, errCode, errMsg, details))
			return
		}
		slog.Error("Failed to create user", "error", err.Error())
		handleServerError(ctx, err)
		return

	}
	ctx.JSON(http.StatusCreated, gin.H{"user_id": userID})
}

// GetUsersReq holds pagination details for retrieving a specified number of users
type GetUsersReq struct {
	PageID   int32 `form:"page_id" binding:"required,gte=1"`
	PageSize int32 `form:"page_size" binding:"required,gte=5,lte=10"`
}

// retrieveUsers route handler for get /users endpoint
func (server *Server) retrieveUsers(ctx *gin.Context) {
	var req GetUsersReq
	if err := bindClientRequest(ctx, &req, querySource); err != nil {
		return
	}
	if !authAccess(ctx, []string{adminRole}) {
		return
	}

	totalUsers, err := server.store.GetTotalNumOfUsers(ctx) //totalUsers
	if err != nil {
		handleServerError(ctx, errors.New("can not retrieve resource from database"))
		return
	}

	users, err := server.store.ListAllUsers(ctx, db.ListAllUsersParams{
		Limit: req.PageSize,
		Offset: req.PageSize * (req.PageID - 1),
	})
	if err != nil {
		handleServerError(ctx, errors.New("can not retrieve resource from database"))
		return
	}

	resp := struct {
		Data []db.ListAllUsersRow `json:"data"`
		Pagination `json:"pagination"`
	}{
		Data: users,
		Pagination: generatePagination(int32(totalUsers), req.PageID, req.PageSize),
	}
	ctx.JSON(http.StatusOK, resp)
}