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
// FIXME: singular noun. //FIXED
// FIXME: present the user with the invalid fields. How can I know what fields should be changed or filled in?
// I knew this was going to come in the review, when I was randomly manually testing the application the last minute before I pushed it :)
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
// FIXME: you should either make params required or provide a default value. They're mutually exclusive.
// I believe it's better to leave the default value here.
func (server *Server) retrieveUsers(ctx *gin.Context) {
	var req GetPaginatedReq
	if err := bindClientRequest(ctx, &req, querySource); err != nil {
		handleClientReqError(ctx, err)
		return
	}
	// FIXME: this logic should be moved in a middleware and abort the request if it doesn't have permissions to perform it. You're bloating the code around.
	// FIXED

	// BUG: you're doing two calls for this operation that can be done at once. You're interacting with the DB one time more than needed.
	// I couldn't find a different way to handle this, but I in the query I am using, I use a cte, to return the total users, and all the users in a byte array
	// I called it TestGetUserPaginated because I wasn't sure it was going to work
	result, err := server.store.TestGetUserPaginated(ctx, db.TestGetUserPaginatedParams{
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
		// BUG: pagination info should be listed first
		// debugged
	}{
		Data:       UserData,
		Pagination: generatePagination(int32(totalUsers), req.PageID, req.PageSize),
	}
	ctx.JSON(http.StatusOK, resp)
}

// getUsersByID route handler for get /users/:user_id, retrieves users by user_id
// FIXME: in the swagger you missed the 401 Unauthorized response
func (server *Server) getUsersByID(ctx *gin.Context) {
	var req UsersIDFromUri
	if err := bindClientRequest(ctx, &req, uriSource); err != nil {
		slog.Error(err.Error())
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
	// FIXME: try to play with the json tag annotation to hide the user password which is the field that led you creating another struct
	// fixed: todo actually, move the db to the internal package
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
// FIXME: should be singular noun this function
// FIXME: the status code should be 202 Accepted, or 200 OK. 204 NoContent is reserved for delete.
func (server *Server) updateUser(ctx *gin.Context) {
	var (
		reqUri  UsersIDFromUri
		reqBody UserReq
	)

	// FIXME: pick the path param & set it to the "id" field in the request payload so the user cannot tamper the data manually.
	// fixed: the user id is not in the req payload
	// Again, use the things Gin provides you. Be more specific with the error handling.
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
	// FIXME: you're doing unnecessary operations in the DB. Delete the user right away. You can decide if trigger an error for non existing user or report success. Be gentle with the DB load.
	// fixed
	numUsersDeleted, err := server.store.DeleteUser(ctx, req.UserID)
	if err != nil {
		slog.Error(err.Error())
		handleServerError(ctx, err)
		return
	}
	if numUsersDeleted < 1 {
		// check if user to be deleted was a master admin, to correctly handle errors
		user_id, err := server.store.GetMasterUser(ctx)
		if err != nil {
			handleServerError(ctx, err)
			return
		}
		if user_id == req.UserID {
			err := errDeleteMaster
			detail := err.Error()
			handleDbError(ctx, err, detail)
			return
		}
		detail := fmt.Sprintf("user with id: %v, doesn't exist", req.UserID)
		err = errNotFound
		handleDbError(ctx, err, detail)
		return
	}
	// FIXME: successMessage() could be omitted
	// fixed
	ctx.Status(http.StatusNoContent)
}

type UpdatePassowrdReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// updateUserPassword route handler for patch /users/:user_id/password/change
// updates a user password
// FIXME: this route could be "/password/change" if we have the "/password/reset"
// I've chosen to just use /password, and remove the reset route, I said something about it in the next fixme tag
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

// FIXME: reset to what? This can be omitted or merged with the "updateUserPassword".
// fixed: using one handler for password update
// Usually, it's a POST that sends a mail asking to reset the password and then you can invoke the above-mentioned PUT.
// I tried sending emails using googles smtp servers, but that service needed some authorization and the company wasn't going to help me get that.
// so the reset password was for the admin to help users restore their passwords to default passwords when they forgot it.
