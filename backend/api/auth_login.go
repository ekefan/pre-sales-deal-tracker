package api

import (
	"errors"
	"net/http"
	"time"

	// "time"

	db "github.com/ekefan/pre-sales-deal-tracker/backend/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// LogingReq holds fields required to authenticate and log in users
type LoginReq struct {
	Username string `json:"username" bindiing:"required,alphanum,gte=4,lte=6"`
	Password string `json:"password" binding:"required,alphanum"`
}

// LoginResp holds the fields in the response body if login is successful
type LoginResp struct {
	AccessToken string        `json:"access_token"`
	UserData    UserLoginResp `json:"user_data"`
}

// authLogin handles client log in
func (server *Server) authLogin(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err, "BAD_REQUEST"))
		return
	}
	// TODO:
	// Update db to set updated_at for userlogin to nullable
	// 3. create json web tokens
	// 4. add token maker to server struct
	// 5. create token

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			ctx.JSON(http.StatusNotFound, errorResponse(err, "NOT_FOUND"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "SERVER_ERROR"))
		return
	}
	if err := ValidatePassword(user.Password, req.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err, "INVALID_PASSWORD"))
		return
	}
	userData := generateUserData(user)
	resp := LoginResp{
		AccessToken: "thisisatestaccesstoken",
		UserData:    userData,
	}

	ctx.JSON(http.StatusOK, resp)
}

// generateUserData generates a userLoginResp struct 
// updatedAt field will be nil if user has never been updated
func generateUserData(user db.User) UserLoginResp {

	if !user.UpdatedAt.Time.IsZero() {
		return UserLoginResp{
			UserID:          user.ID,
			Username:        user.Username,
			Fullname:        user.FullName,
			Role:            user.Role,
			Email:           user.Email,
			PasswordChanged: user.PasswordChanged,
			UpdatedAt:       nil,
		}
	}

	updatedAt := user.UpdatedAt.Time.Format(time.RFC3339)

	return UserLoginResp{
		UserID:          user.ID,
		Username:        user.Username,
		Fullname:        user.FullName,
		Role:            user.Role,
		Email:           user.Email,
		PasswordChanged: user.PasswordChanged,
		UpdatedAt:       &updatedAt,
	}
}
