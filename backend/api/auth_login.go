package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// LogingReq holds fields required to authenticate and log in users
type LoginReq struct {
	Username string `json:"username" binding:"required,alphanum,min=4,max=6"`
	Password string `json:"password" binding:"required"`
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
	AccessToken string `json:"access_token"`
	// T0DO: are we sure we need to send back also this data when the user sign in? // DONE
	// Yes, I used the data to populate the user_profile on the dashboard on the UI
	UserData UserLoginResp `json:"user_data"`
}

// authLogin handles client log in
func (server *Server) authLogin(ctx *gin.Context) {
	var req LoginReq
	// Fixme: it seems overcomplicated this function. //DONE: simplified the function
	// You're using the Gin Web framework, you should use its tools. Try to simplify the code
	if err := bindClientRequest(ctx, &req, jsonSource); err != nil {
		handleClientReqError(ctx, err)
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	// FIXME: you should use errors.Is or errors.As to type check your error.
	if err != nil {
		details := fmt.Sprintf("user with username: %v, doesn't exist", req.Username)
		if handleDbError(ctx, err, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	if err := ValidatePassword(user.Password, req.Password); err != nil {
		handlePasswordValidationError(ctx, err)
		return
	}
	accessToken, err := server.tokenGenerator.GenerateToken(user.ID, user.Role, server.config.TokenDuration)
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
