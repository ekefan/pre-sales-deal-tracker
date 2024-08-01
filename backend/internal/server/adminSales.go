package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	// db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

// LoginReq holds fields required to access user details
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	ID        int64     `json:"user_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	FullName  string    `json:"fullname"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// userLogin handler for loging users to application
func (s *Server) userLogin(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.Store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//TODO: check passwords match
	//Create token or cookies and add to resp
	resp := LoginResp{
		//accessToken:
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		FullName:  user.FullName,
		Email:     user.Email,
		UpdatedAt: user.UpdatedAt.Time,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, resp)

}

type UpdatePitchReq struct {
	ID              int64  `json:"pitch_request_id" binding:"required"`
	Status          string `json:"status" binding:"required"`
	PitchTag        string `json:"pitch_tag" binding:"required"`
	CustomerRequest string `json:"customer_request" binding:"required"`
	AdminViewed     bool   `json:"admin_viewed" binding:"required"`
}

type PitchResp struct {
	ID              int64     `json:"pitch_request_id"`
	SalesRepID      int64     `json:"sales_rep_id"`
	SalesRepName    string    `json:"sales_rep_name"`
	Status          string    `json:"status"`
	CustomerName    string    `json:"customer_name"`
	PitchTag        string    `json:"pitch_tag"`
	CustomerRequest string    `json:"customer_request"`
	RequestDeadline time.Time `json:"request_deadline"`
	AdminViewed     bool      `json:"admin_viewed"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (s *Server) updatePitchReqHandler(ctx *gin.Context) {
	var req UpdatePitchReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// getAuthPayload and require role == admin or sales
	pitchReq, err := s.Store.GetPitchRequestForUpdate(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("pitch request not found:%v", err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	updateTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	args := db.UpdatePitchRequestParams{
		ID:              pitchReq.ID,
		Status:          req.Status,
		PitchTag:        req.PitchTag,
		CustomerRequest: req.CustomerRequest,
		AdminViewed:     req.AdminViewed,
		UpdatedAt:       updateTime,
	}
	updatedPitchReq, err := s.Store.UpdatePitchRequest(ctx, args)
	if err != nil {
		//check possible pq err
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := PitchResp{
		ID:              updatedPitchReq.ID,
		SalesRepID:      updatedPitchReq.SalesRepID,
		SalesRepName:    updatedPitchReq.SalesRepName,
		Status:          updatedPitchReq.Status,
		CustomerName:    updatedPitchReq.CustomerName,
		PitchTag:        updatedPitchReq.PitchTag,
		CustomerRequest: updatedPitchReq.CustomerRequest,
		RequestDeadline: updatedPitchReq.RequestDeadline,
		AdminViewed:     updatedPitchReq.AdminViewed,
		CreatedAt:       updatedPitchReq.CreatedAt,
		UpdatedAt:       updatedPitchReq.UpdatedAt.Time,
	}
	ctx.JSON(http.StatusOK, resp)
}

type ListDealsReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// getDealsHandler http handler for the api end point for getting list of all deals currently
func (s *Server) getDealsHandler(ctx *gin.Context) {
	var req ListDealsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args := db.AdminViewDealsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	users, err := s.Store.AdminViewDeals(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, users)
}