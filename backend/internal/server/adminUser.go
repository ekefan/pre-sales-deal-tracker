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

// CreateUsrReq holds fields that must be provided by client to create user
type CreateUsrReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Role     string `json:"Role" binding:"required,valid-role"`
	FullName string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// CreateUsrRep holds fields that must be provided to client after creating a user
type CreateUsrResp struct {
	Role      string    `json:"Role"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// adminCreateUserHandler http handler for the api end point for creating a new user
func (s *Server) adminCreateUserHandler(ctx *gin.Context) {
	var req CreateUsrReq
	//validate and bind request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Hash Password to prevent saving user password in database
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	args := db.CreateNewUserParams{
		Username:        req.Username,
		Role:            req.Role,
		FullName:        req.FullName,
		Email:           req.Email,
		Password:        passwordHash,
	}

	user, err := s.Store.CreateNewUser(ctx, args)
	if err != nil {
		if pqErrHandler(ctx, "user", err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := CreateUsrResp{
		Role:      user.Role,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, resp)
}

// AdminUpdateUsrReq holds the field - ID to unmarshall json requests
// ID here is the id of the user to be updated
// They are all required, however if no new values are passed... the current
// the current user fields will be used
type AdminUpdateUsrReq struct {
	ID        int64  `json:"user_id" binding:"required"`
	Fullname  string `json:"fullname" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,alphanum"`
}

// AdminUpdateUsrResp holds the fields for responding accurately to updating user end-point
type AdminUpdateUsrResp struct {
	UserID          int64     `json:"user_id"`
	Username        string    `json:"username"`
	Role            string    `json:"role"`
	Fullname        string    `json:"fullname"`
	Email           string    `json:"email"`
	PasswordChanged bool      `json:"password_changed"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// adminUpdateUserHandler http handler for the api end point for updating a user
func (s *Server) adminUpdateUserHandler(ctx *gin.Context) {
	var req AdminUpdateUsrReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//get access token
	if !authAccess(ctx, utils.AdminRole) {
		return
	}
	usr, err := s.Store.GetUserForUpdate(ctx, req.ID)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// Set update time to time now....
	updateTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	// Hash password
	args := db.AdminUpdateUserParams{
		ID:        usr.ID,
		FullName:  req.Fullname,
		Email:     req.Email,
		Username:  req.Username,
		UpdatedAt: updateTime,
	}
	// get
	newUsr, err := s.Store.AdminUpdateUser(ctx, args)
	if err != nil {
		if pqErrHandler(ctx, "user", err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := AdminUpdateUsrResp{
		UserID:          newUsr.ID,
		Username:        newUsr.Username,
		Role:            newUsr.Role,
		Fullname:        newUsr.FullName,
		Email:           newUsr.Email,
		UpdatedAt:       newUsr.UpdatedAt.Time,
		PasswordChanged: newUsr.PasswordChanged,
		CreatedAt:       newUsr.CreatedAt,
	}
	ctx.JSON(http.StatusOK, resp)
}

// AdminDeleteUserReq holds field user id that is to be deleted
type AdminDeleteUserReq struct {
	ID        int64  `uri:"id" binding:"required"`
	// AdminRole string `uri:"admin_role" binding:"required"`
}

// adminDeleteUserhandler http handler for the api end point for Deleting a user
func (s *Server) adminDeleteUserHandler(ctx *gin.Context) {
	var req AdminDeleteUserReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	// authenticated access
	if !authAccess(ctx, utils.AdminRole) {
		return
	}
	exists, err := s.Store.AdminUserExists(ctx, req.ID)
	if err != nil || !exists {
		ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user doesn't exist")))
		return
	}

	err = s.Store.AdminDeleteUser(ctx, req.ID)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successful",
	})
}
