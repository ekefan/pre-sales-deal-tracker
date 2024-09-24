package server

import (
	"fmt"
	"net/http"
	"time"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/ekefan/deal-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

// this handler would be called after a use logs in and their password_changed status is false
type UpdatePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	Password    string `json:"password_update" binding:"required,min=6"`
	UserID      int64  `json:"user_id"`
}

// updatePassWordLoggedIn takes the current password then updates the user password
// FIXME: try mergin these endpoint with the password reset one. Maybe, you can do the following: keep a single endpoint that provide users with the ability to change the password. Have something that mark as expired passwords older than x days. We can discuss this as well.
// DONE: Yes let's discuss it... Just for some context, the reset password is an admin authorized endpoint and it's for when users forget their own password, not that it can't be merged with this but I just did it that way... when I discuss with you I will have a different view. :)
func (s *Server) updatePassWordLoggedIn(ctx *gin.Context) {
	var req UpdatePasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// maybe using session management.. check if user is loggedin
	if !multipleAuthAccess(ctx, []string{utils.AdminRole, utils.ManagerRole, utils.SalesRole}) {
		return
	}
	usr, err := s.Store.GetUserForUpdate(ctx, req.UserID)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	if !utils.CheckPasswordHash(req.OldPassword, usr.Password) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("password invalid: %v", err)))
		return
	}
	success := s.Store.UpdatePassWord(ctx, db.UpdatePassWordParams{
		ID:              req.UserID,
		Password:        req.Password,
		PasswordChanged: true,
		UpdatedAt:       time.Now(),
	})
	if success != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{
		"mesaage": "succesful",
	})
}

// resetPasswordReq holds the id of the user whose password should be updated
type ResetPasswordReq struct {
	UserToUpdateID int64 `json:"user_id" binding:"required,gt=0"`
}

// resetPassword: endpoint for resetting user password to default password
// sends user id from reset password request to perform password reset
// transaction returns successful if successful
// FIXME: the password should be a subresource of the user resource.
// since it's not common to handle the credential in your web server, try do come up with a valuable solution and let's have a discussion around.
// DONE: Yes please, I am anticipating that discussion
func (s *Server) resetPassword(ctx *gin.Context) {
	var req ResetPasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check if user is authorized
	if !authAccess(ctx, utils.AdminRole) {
		return
	}
	success := s.Store.ResetPasswordTxn(ctx, db.ResetPasswordArgs{
		UserID: req.UserToUpdateID,
	})

	if success != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(success))
		return
	}

	// send user the newPassword to their email....
	ctx.JSON(http.StatusNoContent, gin.H{
		"mesaage": "succesful",
	})
}