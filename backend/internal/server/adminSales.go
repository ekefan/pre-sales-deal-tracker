package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	// db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/ekefan/deal-tracker/internal/utils"
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

	// check if the user has been created
	user, err := s.Store.GetUser(ctx, req.Username)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// if user has not updated password...
	// redirect user to update password...
	//set password,
	//update newPassword
	//
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("password invalid: %v", err)))
		return
	}

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
	ID              int64    `json:"pitch_request_id" binding:"required"`
	Status          string   `json:"status" binding:"required"`
	PitchTag        string   `json:"pitch_tag" binding:"required"`
	CustomerRequest string   `json:"customer_request" binding:"required"`
	AdminViewed     bool     `json:"admin_viewed"`
	RequestDealine  UnixTime `json:"request_deadline"`
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
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	updateTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	var deadline time.Time
	if req.RequestDealine.Valid {
		deadline = req.RequestDealine.Time
	} else {
		deadline = pitchReq.RequestDeadline
	}
	args := db.UpdatePitchRequestParams{
		ID:              pitchReq.ID,
		Status:          req.Status,
		PitchTag:        req.PitchTag,
		CustomerRequest: req.CustomerRequest,
		AdminViewed:     req.AdminViewed,
		UpdatedAt:       updateTime,
		RequestDeadline: deadline,
	}
	updatedPitchReq, err := s.Store.UpdatePitchRequest(ctx, args)
	if err != nil {
		if pqErrHandler(ctx, "pitch_requests", err) {
			return
		}
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

// type ListDealsReq struct {
// 	PageID   int32 `form:"page_id" binding:"required,min=1"`
// 	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
// }

// // getDealsHandler http handler for the api end point for getting list of all deals currently
// func (s *Server) getDealsHandler(ctx *gin.Context) {
// 	var req ListDealsReq
// 	if err := ctx.ShouldBindQuery(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// 	args := db.AdminViewDealsParams{
// 		Limit:  req.PageSize,
// 		Offset: (req.PageID - 1) * req.PageSize,
// 	}
// 	deals, err := s.Store.AdminViewDeals(ctx, args)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, deals)
// }

type OngoingDealsReq struct {
	Status string `json:"status" binding:"required"`
}

func (s *Server) getOngoingDeals(ctx *gin.Context) {
	var req OngoingDealsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deals, err := s.Store.GetDealsByStatus(ctx, req.Status)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, deals)

}

type FilterDealReq struct {
	CustomerName    string         `json:"customer_name"`
	ServiceToRender string         `json:"service_to_render"`
	Status          string         `json:"status"`
	MaxProfit       sql.NullString `json:"max_profit"`
	MinProfit       sql.NullString `json:"min_profit"`
	Awarded         bool           `json:"awarded"`
	SalesRepName    string         `json:"sales_rep_name"`
	PageSize        int32          `json:"page_size"`
	PageID          int32          `json:"page_id"`
}

func (s *Server) getFilteredDeals(ctx *gin.Context) {
	var req FilterDealReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.FilterDealsParams{
		CustomerName: req.CustomerName,
		ServiceToRender: req.ServiceToRender,
		Status: req.Status,
		Profit: req.MaxProfit,
		Profit_2: req.MinProfit,
		Awarded: req.Awarded,
		SalesRepName: req.SalesRepName,
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	deals, err := s.Store.FilterDeals(ctx, args)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, deals)

}
