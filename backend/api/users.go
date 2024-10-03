package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	db "github.com/ekefan/pre-sales-deal-tracker/backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

// CreateUserReq holds fields needed to create a user resource
type CreateUserReq struct {
	Username string `json:"username" binding:"required,gte=4,lte=6,alphanum"`
	Fullname string `json:"fullname" binding:"required,gte=4"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=admin sales manager"`
}

// createUser route handler post /users, creates users resource
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
		errMsg = "can not create user"
		details = fmt.Sprintf("user exists with %v or %v", req.Email, req.Username)
		if pgxError(ctx, err, errMsg, details) {
			slog.Warn("user already exist", "username", req.Username, "email", req.Email)
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

// retrieveUsers route handler for get /users, retrieves all users
func (server *Server) retrieveUsers(ctx *gin.Context) {
	var req GetUsersReq
	if err := bindClientRequest(ctx, &req, querySource); err != nil {
		slog.Error(err.Error())
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
		Limit:  req.PageSize,
		Offset: req.PageSize * (req.PageID - 1),
	})
	if err != nil {
		handleServerError(ctx, errors.New("can not retrieve resource from database"))
		return
	}

	resp := struct {
		Data       []db.ListAllUsersRow `json:"data"`
		Pagination `json:"pagination"`
	}{
		Data:       users,
		Pagination: generatePagination(int32(totalUsers), req.PageID, req.PageSize),
	}
	ctx.JSON(http.StatusOK, resp)
}

// UsersIDFromUri holds the uri field user_id
type UsersIDFromUri struct {
	UserID int64 `uri:"user_id"`
}


// getUsersByID route handler for get /users/:user_id, retrieves users by user_id
func (server *Server) getUsersByID(ctx *gin.Context) {
	var req UsersIDFromUri
	if err := bindClientRequest(ctx, &req, uriSource); err != nil {
		slog.Error(err.Error())
		return
	}
	if !authAccess(ctx, []string{adminRole, salesRole, managerRole}) {	
		return
	}
	user, err := server.store.GetUserByID(ctx, req.UserID)
	if err != nil {
		errMsg = "user not found"
		details = fmt.Sprintf("user with user_id: %v, not found", req.UserID)
		if pgxError(ctx, err, errMsg, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	resp := db.ListAllUsersRow{
		UserID: user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Role: user.Role,
		Email: user.Email,
		PasswordChanged: user.PasswordChanged,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Time.Format(time.RFC3339),
	}
	ctx.JSON(http.StatusOK, resp)

}