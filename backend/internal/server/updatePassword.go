package server

import (
	"fmt"
	"net/http"
	"time"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/ekefan/deal-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

// / this handler would be called after a use logs in and their password_changed status is false
type UpdatePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	Password    string `json:"password_update" binding:"required,min=6"`
	UserID      int64  `json:"user_id"`
}

// updatePassWordLoggedIn takes the current password then updates the user password
func (s *Server) updatePassWordLoggedIn(ctx *gin.Context) {
	var req UpdatePasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//maybe using session management.. check if user is loggedin
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
	ctx.JSON(http.StatusOK, gin.H{
		"mesaage": "succesful",
	})
}

// resetPasswordReq holds the id of the user whose password should be updated
type resetPasswordReq struct {
	UserToUpdateID int64 `json:"user_id" binding:"required,gt=0"`
}


// resetPassword: endpoint for resetting user password to default password
// sends user id from reset password request to perform password reset 
// transaction returns successful if successful
func (s *Server) resetPassword(ctx *gin.Context) {
	var req resetPasswordReq
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

	//send user the newPassword to their email....
	ctx.JSON(http.StatusOK, gin.H{
		"mesaage": "succesful",
	})
}