package server

import (
	"database/sql"
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
	if !authAccess(ctx, utils.AdminRole) || !authAccess(ctx, utils.ManagerRole) || !authAccess(ctx, utils.SalesRole) {
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
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
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


//UnderConstruction
// type ForgotPasswordReq struct {
// 	Email string `json:"email" binding:"required,email"`
// }

// func (s *Server) forgotPassword(ctx *gin.Context) {
// 	var req ForgotPasswordReq
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	usr, err := s.Store.ForgotPassword(ctx, req.Email)
// 	if err != nil {
// 		if sqlNoRowsHandler(ctx, err) {
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	newRandomPass := randomPasswordCode()
// 	updatePasswordArgs := db.UpdatePassWordParams{
// 		ID:              usr.ID,
// 		Password:        newRandomPass,
// 		PasswordChanged: true,
// 		UpdatedAt: sql.NullTime{
// 			Time:  time.Now(),
// 			Valid: true,
// 		},
// 	}
// 	success := s.Store.UpdatePassWord(ctx, updatePasswordArgs)
// 	if success != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	//send user the newPassword to their email....
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"mesaage": "succesful",
// 	})
// }