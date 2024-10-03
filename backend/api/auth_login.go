package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// LogingReq holds fields required to authenticate and log in users
type LoginReq struct {
	Username string `json:"username" binding:"required,alphanum,gte=4,lte=6"`
	Password string `json:"password" binding:"required,alphanum"`
}

// UserLoginResp sub field in the login response body
type UserLoginResp struct {
	UserID          int64  `json:"user_id"`
	Username        string `json:"username"`
	FullName        string `json:"full_name"`
	Role            string `json:"role"`
	Email           string `json:"email"`
	PasswordChanged bool   `json:"password_changed"`
	UpdatedAt       string `json:"updated_at"`
	CreatedAt       string `json:"created_at"`
}

// LoginResp holds the fields in the response body if login is successful
type LoginResp struct {
	AccessToken string        `json:"access_token"`
	UserData    UserLoginResp `json:"user_data"`
}

// authLogin handles client log in
func (server *Server) authLogin(ctx *gin.Context) {
	var req LoginReq
	if err := bindClientRequest(ctx, &req, jsonSource); err != nil {
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		errMsg = "user not found"
		details = fmt.Sprintf("no user with username %v exists", req.Username)
		if pgxError(ctx, err, errMsg, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	if err := ValidatePassword(user.Password, req.Password); err != nil {
		statusCode = http.StatusUnauthorized
		errCode = "UNAUTHORIZED"
		errMsg = "invalid password"
		details = "the user password sent by this client is incorrect"
		ctx.JSON(statusCode, errorResponse(statusCode, errCode, errMsg, details))
		return
	}
	accessToken, _, err := server.tokenGenerator.GenerateToken(user.ID, user.Role, server.config.TokenDuration)
	if err != nil {
		handleServerError(ctx, err)
		return
	}
	userData := UserLoginResp{
		UserID:          user.ID,
		Username:        user.Username,
		FullName:        user.FullName,
		Role:            user.Role,
		Email:           user.Email,
		PasswordChanged: user.PasswordChanged,
		UpdatedAt:       user.UpdatedAt.Time.Format(time.RFC3339),
		CreatedAt:       user.CreatedAt.Time.Format(time.RFC3339),
	}

	resp := LoginResp{
		AccessToken: accessToken,
		UserData:    userData,
	}

	ctx.JSON(http.StatusOK, resp)
}
