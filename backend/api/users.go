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

// UserReq holds fields needed to create or update a user resource
type UserReq struct {
	Username string `json:"username" binding:"required,gte=4,lte=6,alphanum"`
	FullName string `json:"full_name" binding:"required,gte=4"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=admin sales manager"`
}

// createUser route handler post /users, creates users resource
func (server *Server) createUsers(ctx *gin.Context) {
	var req UserReq
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
		FullName: req.FullName,
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
	UserID int64 `uri:"user_id" binding:"required"`
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
		UserID:          user.ID,
		Username:        user.Username,
		FullName:        user.FullName,
		Role:            user.Role,
		Email:           user.Email,
		PasswordChanged: user.PasswordChanged,
		CreatedAt:       user.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:       user.UpdatedAt.Time.Format(time.RFC3339),
	}
	ctx.JSON(http.StatusOK, resp)

}

// updateUsers route handler for put /users/:user_id
func (server *Server) updateUsers(ctx *gin.Context) {
	var (
		reqUri  UsersIDFromUri
		reqBody UserReq
	)

	uriErr := bindClientRequest(ctx, &reqUri, uriSource)
	reqBodyErr := bindClientRequest(ctx, &reqBody, jsonSource)
	if uriErr != nil || reqBodyErr != nil {
		slog.Error("failed to bind client request", "uri error", uriErr, "req body err", reqBodyErr)
		return
	}

	if !authAccess(ctx, []string{adminRole}) {
		return
	}

	usr, err := server.store.GetUserByID(ctx, reqUri.UserID)
	if err != nil {
		errMsg = "user not found"
		details = fmt.Sprintf("user with user_id: %v, not found", reqUri.UserID)
		if pgxError(ctx, err, errMsg, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	args := db.UpdateUserParams{
		ID: reqUri.UserID,
		Username: reqBody.Username,
		FullName: reqBody.FullName,
		Role:     reqBody.Role,
		Email:    reqBody.Email,
	}
	err = server.store.UpdateUserTx(ctx, db.UpdateUserTxParams{
		UpdateUserParams: args,
		OldFullName: usr.FullName,
	})
	if err != nil {
		errMsg = "can not update user"
		details = fmt.Sprintf("user with username: %v, or email: %v, exists", reqBody.Username, reqBody.Email)
		if pgxError(ctx, err, errMsg, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, successMessage())
}

// deleteUsers route handler for delete /users/:user_id, deletes user with user_id
func (server *Server) deleteUsers(ctx *gin.Context) {
	var req UsersIDFromUri
	if err := bindClientRequest(ctx, &req, uriSource); err != nil {
		slog.Error(err.Error())
		return
	}

	if !authAccess(ctx, []string{adminRole}) {
		return
	}
	if _, err := server.store.GetUserByID(ctx, req.UserID); err != nil {
		errMsg = "user not found"
		details = fmt.Sprintf("user with user_id: %v, not found", req.UserID)
		if pgxError(ctx, err, errMsg, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	if err := server.store.DeleteUser(ctx, req.UserID); err != nil {
		handleServerError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, successMessage())
}

type UpdatePassowrdReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// updateUserPassword route handler for patch /users/:user_id/password
// updates a user password
func (server *Server) updateUserPassword(ctx *gin.Context) {
	var (
		reqUri  UsersIDFromUri
		reqBody UpdatePassowrdReq
	)

	uriErr := bindClientRequest(ctx, &reqUri, uriSource)
	reqBodyErr := bindClientRequest(ctx, &reqBody, jsonSource)
	if uriErr != nil || reqBodyErr != nil {
		slog.Error("failed to bind client request", "uri error", uriErr, "req body err", reqBodyErr)
		return
	}

	if !authAccess(ctx, []string{adminRole, managerRole, salesRole}) {
		return
	}

	user, err := server.store.GetUserByID(ctx, reqUri.UserID)
	if err != nil {
		errMsg = "user not found"
		details = fmt.Sprintf("user with user_id: %v, not found", reqUri.UserID)
		if pgxError(ctx, err, errMsg, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	if err := ValidatePassword(user.Password, reqBody.OldPassword); err != nil {
		statusCode = http.StatusUnauthorized
		errCode = "UNAUTHORIZED"
		errMsg = "invalid password"
		details = fmt.Sprintf("the old password sent doesn't not match with password of user with id: %v", reqUri.UserID)
		ctx.JSON(statusCode, errorResponse(statusCode, errCode, errMsg, details))
		return
	}

	hash, err := HashPassword(reqBody.NewPassword)
	if err != nil {
		handleServerError(ctx, err)
		return
	}
	
	if err := server.store.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:              reqUri.UserID,
		Password:        hash,
		PasswordChanged: true,
	}); err != nil {
		handleServerError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

// resetUserPassword route handler for patch /users/:user_id/password
// resets a user password
func (server *Server) resetUserPassword(ctx *gin.Context) {
	var reqUri UsersIDFromUri

	uriErr := bindClientRequest(ctx, &reqUri, uriSource)
	if uriErr != nil {
		slog.Error("failed to bind client request", "uri error", uriErr.Error())
		return
	}

	if !authAccess(ctx, []string{adminRole}) {
		return
	}

	if _, err := server.store.GetUserByID(ctx, reqUri.UserID); err != nil {
		errMsg = "user not found"
		details = fmt.Sprintf("user with user_id: %v, not found", reqUri.UserID)
		if pgxError(ctx, err, errMsg, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	hash, err := HashPassword(db.DefaultUserPassword)
	if err != nil {
		handleServerError(ctx, err)
		return
	}
	if err := server.store.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:              reqUri.UserID,
		Password:        hash,
		PasswordChanged: false,
	}); err != nil {
		handleServerError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}
