package api

import (
	"errors"
	"log/slog"
	"net/http"

	db "github.com/ekefan/pre-sales-deal-tracker/backend/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)
// CreateUserReq holds fields needed to create a user resource
type CreateUserReq struct {
	Username string `json:"username" binding:"required,gte=4,lte=6,alphanum"`
	Fullname string `json:"fullname" binding:"required,gte=4"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=admin sales manager"`
}

// createUser api handler for creating user resource
func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserReq
	if err := bindClientRequest(ctx, &req, jsonSource); err != nil {
		slog.Error(err.Error())
		return
	}
	if !authAccess(ctx, []string{adminRole}) {
		return
	}
	hash, err := HashPassword(db.DefaultUserPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "SEVER ERROR"))
		return
	}
	userID, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		FullName: req.Fullname,
		Email:    req.Email,
		Role:     req.Role,
		Password: hash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			slog.Warn("user already exist", "username", req.Username, "email", req.Email)
			ctx.JSON(http.StatusConflict, errorResponse(err, "USER_ALREADY_EXISTS"))
			return
		}
		slog.Error("Failed to create user", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "SERVER_ERROR"))
		return
		
	}
	ctx.JSON(http.StatusCreated, gin.H{"user_id": userID})
}
