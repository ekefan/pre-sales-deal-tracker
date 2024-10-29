package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	db "github.com/ekefan/pre-sales-deal-tracker/backend/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

// UserReq holds fields needed to create or update a user resource
type UserReq struct {
	Username string `json:"username" binding:"required,min=4,max=6,alphanum"`
	FullName string `json:"full_name" binding:"required,min=4"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=admin sales"`
}

// createUser route handler post /users, creates a user
func (server *Server) createUser(ctx *gin.Context) {
	var req UserReq
	if err := bindClientRequest(ctx, &req, jsonSource); err != nil {
		handleClientReqError(ctx, err)
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
		details := fmt.Sprintf("user exists with %v or %v", req.Email, req.Username)
		if handleDbError(ctx, err, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"user_id": userID})
}

// retrieveUsers route handler for get /users, retrieves all users
func (server *Server) retrieveUsers(ctx *gin.Context) {
	var req GetPaginatedReq
	if err := bindClientRequest(ctx, &req, querySource); err != nil {
		handleClientReqError(ctx, err)
		return
	}
	result, err := server.store.GetUserPaginated(ctx, db.GetUserPaginatedParams{
		Limit:  req.PageSize,
		Offset: req.PageSize * (req.PageID - 1),
	})
	if err != nil {
		slog.Error(err.Error())
		handleServerError(ctx, err)
		return
	}
	UserData := []User{}
	totalUsers := result.TotalUsers
	if len(result.Users) > 0 {
		usRrr := json.Unmarshal(result.Users, &UserData)
		if usRrr != nil {
			slog.Error(usRrr.Error())
			handleServerError(ctx, err)
			return
		}
	}
	resp := struct {
		Pagination `json:"pagination"`
		Data       []User `json:"data"`

	}{
		Pagination: generatePagination(int32(totalUsers), req.PageID, req.PageSize),
		Data:       UserData,
	}
	ctx.JSON(http.StatusOK, resp)
}

// getUsersByID route handler for get /users/:user_id, retrieves users by user_id
func (server *Server) getUsersByID(ctx *gin.Context) {
	var req UsersIDFromUri
	if err := bindClientRequest(ctx, &req, uriSource); err != nil {
		handleClientReqError(ctx, err)
		return
	}
	user, err := server.store.GetUserByID(ctx, req.UserID)
	if err != nil {
		details := fmt.Sprintf("user with user_id: %v, not found", req.UserID)
		if handleDbError(ctx, err, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	// using sqlc, I couldn't find a way to make exclude a field using json tags. I created api models for that instead...
	resp := User{
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
func (server *Server) updateUser(ctx *gin.Context) {
	var (
		reqUri  UsersIDFromUri
		reqBody UserReq
	)
	uriErr := bindClientRequest(ctx, &reqUri, uriSource)
	if uriErr != nil {
		handleClientReqError(ctx, uriErr)
		return
	}
	reqBodyErr := bindClientRequest(ctx, &reqBody, jsonSource)
	if reqBodyErr != nil {
		handleClientReqError(ctx, reqBodyErr)
		return
	}

	usr, err := server.store.GetUserByID(ctx, reqUri.UserID)
	if err != nil {
		details := fmt.Sprintf("user with user_id: %v, not found", reqUri.UserID)
		if handleDbError(ctx, err, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	args := db.UpdateUserParams{
		ID:       reqUri.UserID,
		Username: reqBody.Username,
		FullName: reqBody.FullName,
		Role:     reqBody.Role,
		Email:    reqBody.Email,
	}
	// [Q]: if we look merely to the user entity, I don't see the reason why we put it in a transaction.
	// [A]: the deals have the name of the sales rep who brought them, so if a user being updated is in sales, the deal get's updated to
	err = server.store.UpdateUserTx(ctx, db.UpdateUserTxParams{
		UpdateUserParams: args,
		OldFullName:      usr.FullName,
	})
	if err != nil {
		details := fmt.Sprintf("user with username: %v, or email: %v, exists", reqBody.Username, reqBody.Email)
		if handleDbError(ctx, err, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, successMessage())
}

// UsersIDFromUri holds the uri field user_id
type UsersIDFromUri struct {
	UserID int64 `uri:"user_id" binding:"required,min=1,numeric"`
}

// deleteUsers route handler for delete /users/:user_id, deletes user with user_id
func (server *Server) deleteUser(ctx *gin.Context) {
	var req UsersIDFromUri
	if err := bindClientRequest(ctx, &req, uriSource); err != nil {
		handleClientReqError(ctx, err)
		return
	}
	err := server.store.StoreDeleteUser(ctx, req.UserID)
	if err != nil {
		if err.Error() == errNotFound.Error() ||
			err.Error() == errDeleteMaster.Error() {
			handleDbError(ctx, err, err.Error())
			return
		}
		slog.Error(err.Error())
		handleServerError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

type UpdatePassowrdReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// updateUserPassword route handler for patch /users/:user_id/password/change
// updates a user password
func (server *Server) updateUserPassword(ctx *gin.Context) {
	var (
		reqUri  UsersIDFromUri
		reqBody UpdatePassowrdReq
	)

	uriErr := bindClientRequest(ctx, &reqUri, uriSource)
	if uriErr != nil {
		handleClientReqError(ctx, uriErr)
		return
	}
	reqBodyErr := bindClientRequest(ctx, &reqBody, jsonSource)
	if reqBodyErr != nil {
		handleClientReqError(ctx, reqBodyErr)
		return
	}
	user, err := server.store.GetUserByID(ctx, reqUri.UserID)
	if err != nil {
		details := fmt.Sprintf("user with user_id: %v, not found", reqUri.UserID)
		if handleDbError(ctx, err, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	if err := ValidatePassword(user.Password, reqBody.OldPassword); err != nil {
		handlePasswordValidationError(ctx, err)
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
	ctx.Status(http.StatusNoContent)
}