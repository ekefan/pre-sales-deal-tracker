package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// CreateDealReq holds fields needed to create a new deal
type CreateDealReq struct {
	PitchID             int64  `json:"pitch_id" binding:"required"`
	SalesRepName        string `json:"sales_rep_name" binding:"required"`
	CustomerName        string `json:"customer_name" binding:"required"`
	ServiceToRender     string `json:"service_to_render" binding:"required"`
	Status              string `json:"status" binding:"required"`
	StatusTag           string `json:"status_tag" binding:"required"`
	CurrentPitchRequest string `json:"current_pitch_request" binding:"required"`
}

type CreateDealResp struct {
	ID                  int64     `json:"id"`
	PitchID             int64     `json:"pitch_id"`
	SalesRepName        string    `json:"sales_rep_name"`
	CustomerName        string    `json:"customer_name"`
	ServiceToRender     string    `json:"service_to_render"`
	Status              string    `json:"status"`
	StatusTag           string    `json:"status_tag"`
	CurrentPitchRequest string    `json:"current_pitch_request"`
	NetTotalCost        string    `json:"net_total_cost"`
	Profit              string    `json:"profit"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	ClosedAt            time.Time `json:"closed_at"`
	Awarded             bool      `json:"awarded"`
}

// adminCreateDealHandler http handler for the api end point for creating a new deal
func (s *Server) adminCreateDealHandler(ctx *gin.Context) {
	//validate json-req and unmarshall it to the req
	var req CreateDealReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//find a way to perform authorization
	// -- maybe receive user id from token/cookies and verify it's an admin

	//make call to database
	deal, err := s.Store.CreateDeal(ctx, db.CreateDealParams{
		PitchID:             db.SetNullPitchID(req.PitchID),
		SalesRepName:        req.SalesRepName,
		CustomerName:        req.CustomerName,
		ServiceToRender:     req.ServiceToRender,
		Status:              req.Status,
		StatusTag:           req.StatusTag,
		CurrentPitchRequest: req.CurrentPitchRequest,
	})

	if err != nil {
		if _, ok := err.(*pq.Error); ok {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	
	resp := CreateDealResp{
		ID:                  deal.ID,
		PitchID:             deal.PitchID.Int64, // .Int64 returns the value from the sql.NullInt64
		SalesRepName:        deal.SalesRepName,
		CustomerName:        deal.CustomerName,
		ServiceToRender:     deal.ServiceToRender,
		Status:              deal.Status,
		StatusTag:           deal.StatusTag,
		CurrentPitchRequest: deal.CurrentPitchRequest,
		NetTotalCost:        deal.NetTotalCost.String,
		Profit:              deal.Profit.String,
		CreatedAt:           deal.CreatedAt,
		UpdatedAt:           deal.UpdatedAt.Time,
		ClosedAt:            deal.ClosedAt.Time,
		Awarded:             deal.Awarded,
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateDealReq holds fields used to update a deal
type UpdateDealReq struct {
	ID                  int64    `json:"id" binding:"required"`
	ServiceToRender     string   `json:"service_to_render" binding:"required"`
	Status              string   `json:"status" binding:"required"`
	StatusTag           string   `json:"status_tag" binding:"required"`
	CurrentPitchRequest string   `json:"current_pitch_request" binding:"required"`
	UpdatedAt           UnixTime `json:"updated_at" binding:"required"`
	ClosedAt            UnixTime `json:"closed_at"`
}

// adminUpdateDealsHandler http handler for the api end point for updating a deal
func (s *Server) adminUpdateDealHandler(ctx *gin.Context) {
	var req UpdateDealReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//verify that the resouces is accessed by the admin only, check for role in
	// authorization payload

	deal, err := s.Store.AdminGetDealForUpdate(ctx, req.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			ctx.JSON(http.StatusForbidden, errorResponse(pqErr))
			return
		}
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	updatedAt := sql.NullTime{
		Time:  req.UpdatedAt.Time,
		Valid: req.UpdatedAt.Valid,
	}
	closedAt := sql.NullTime{
		Time:  req.UpdatedAt.Time,
		Valid: req.ClosedAt.Valid,
	}
	updatedDeal, err := s.Store.AdminUpdateDeal(ctx, db.AdminUpdateDealParams{
		ID:                  deal.ID,
		ServiceToRender:     req.ServiceToRender,
		Status:              req.Status,
		StatusTag:           req.StatusTag,
		CurrentPitchRequest: req.CurrentPitchRequest,
		UpdatedAt:           updatedAt,
		ClosedAt:            closedAt,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			//check for specific pq Errors but...
			ctx.JSON(http.StatusForbidden, errorResponse(pqErr))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := CreateDealResp{
		ID:                  updatedDeal.ID,
		PitchID:             updatedDeal.PitchID.Int64,
		SalesRepName:        updatedDeal.SalesRepName,
		CustomerName:        updatedDeal.CustomerName,
		ServiceToRender:     updatedDeal.ServiceToRender,
		Status:              updatedDeal.Status,
		StatusTag:           updatedDeal.StatusTag,
		CurrentPitchRequest: updatedDeal.CurrentPitchRequest,
		NetTotalCost:        updatedDeal.NetTotalCost.String,
		Profit:              updatedDeal.Profit.String,
		CreatedAt:           updatedDeal.CreatedAt,
		UpdatedAt:           updatedDeal.UpdatedAt.Time,
		ClosedAt:            updatedDeal.ClosedAt.Time,
		Awarded:             updatedDeal.Awarded,
	}

	ctx.JSON(http.StatusOK, resp)
}

// AdminDeleteUserhandler holds field user id that is to be deleted
type DeleteDealReq struct {
	ID int64 `uri:"deal_id" binding:"required"`
}

// adminDeleteDealHandler http handler for the api end point for Deleting a deal
func (s *Server) adminDeleteDealHandler(ctx *gin.Context) {
	var req DeleteDealReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//get payload and check the for user.role

	exists, err := s.Store.AdminDealExists(ctx, req.ID)
	if err != nil || !exists {
		ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("deal doesn't exist")))
		return
	}
	err = s.Store.AdminDeleteDeal(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("couldn't delete deal: %v", err)))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successful",
	})
}

type ListUsersReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listUsershandler http handler for the api end point for getting list of users currently
func (s *Server) listUsersHandler(ctx *gin.Context) {
	var req ListUsersReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args := db.AdminViewUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	users, err := s.Store.AdminViewUsers(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// ================TODO=============================
//1. Create Custom validation for role: oneof - admin, sales, manager
//2. Create hash password functionality
//3. send user email to update password, login
//4. Must write test to validate that update deal works well
