package api

import (
	"net/http"
	"time"

	db "github.com/ekefan/pre-sales-deal-tracker/backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

// LogingReq holds fields required to authenticate and log in users
type LoginReq struct {
	Username string `json:"username" bindiing:"required,alphanum,gte=4,lte=6"`
	Password string `json:"password" binding:"required,alphanum"`
}

// LoginResp holds the fields in the response body if login is successful
type LoginResp struct {
	AccessToken string `json:"access_token"`
	UserData UserLoginResp `json:"user_data"`
}

// authLogin handles client log in
func (server *Server) authLogin(ctx *gin.Context) {
	var req LoginReq 
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err, "BAD_REQUEST"))
		return
	}
	// TODO: 
	// 1. create password encrypter and checker
	// 2. Verfiy passwor
	// 3. create json web tokens
	// 4. add token maker to server struct
	// 5. create token
	user, err := server.store.AuthLogin(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err, "NOT_FOUND"))
		return
	}

	userData := generateLoginUserData(user)
	resp := LoginResp{
		AccessToken: "thisisatestaccesstoken",
		UserData: userData,
	}

	ctx.JSON(http.StatusOK, resp)
}


func generateLoginUserData(user db.User) UserLoginResp {
	if !user.PasswordChanged {
		return UserLoginResp{
			UserID: user.ID,
			Username: user.Username,
			Fullname: user.FullName,
			Role: user.Role,
			Email: user.Email,
			PasswordChanged: user.PasswordChanged,
		}
	}
	
	return UserLoginResp{
		UserID: user.ID,
		Username: user.Username,
		Fullname: user.FullName,
		Role: user.Role,
		Email: user.Email,
		PasswordChanged: user.PasswordChanged,
		UpdatedAt: user.UpdatedAt.Time.Format(time.RFC3339),
	}
}