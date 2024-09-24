package server

import (
	"fmt"
	"net/http" 

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/ekefan/deal-tracker/internal/token"
	"github.com/ekefan/deal-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

// PitchReq holds fields for creating a Pitch Request
type PitchReq struct {
	SalesRepID      int64    `json:"sales_rep_id" binding:"required"`
	SalesRepName    string   `json:"sales_rep_name" binding:"required"`
	Status          string   `json:"status" binding:"required"`
	CustomerName    string   `json:"customer_name" binding:"required"`
	PitchTag        string   `json:"pitch_tag" binding:"required"`
	CustomerRequests []string   `json:"customer_requests" binding:"required"`
	RequestDeadline UnixTime `json:"request_deadline" binding:"required"`
}

// salesCreatePitchReqHandler api endpoint for creating a pitch request
// expected to be called by a user with sales role
// To be successful a PitchReq is sent to the handler
func (s *Server) salesCreatePitchReqHandler(ctx *gin.Context) {
	var req PitchReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	/// General Problem of Parsing time between time.Time and json

	// instead of receiving salesRepID from json validate it through payload
	if !authAccess(ctx, utils.SalesRole) {
		return
	}
	args := db.CreatePitchRequestParams{
		SalesRepID:      req.SalesRepID,
		SalesRepName:    req.SalesRepName,
		Status:          req.Status,
		CustomerName:    req.CustomerName,
		PitchTag:        req.PitchTag,
		CustomerRequest: req.CustomerRequests,
		RequestDeadline: req.RequestDeadline.Time,
	}
	pitchRequest, err := s.Store.CreatePitchRequest(ctx, args)
	if err != nil {
		//check possible pq.Error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := PitchResp{
		ID:              pitchRequest.ID,
		SalesRepID:      pitchRequest.SalesRepID,
		SalesRepName:    pitchRequest.SalesRepName,
		Status:          pitchRequest.Status,
		CustomerName:    pitchRequest.CustomerName,
		PitchTag:        pitchRequest.PitchTag,
		CustomerRequests: pitchRequest.CustomerRequest,
		RequestDeadline: pitchRequest.RequestDeadline.Unix(),
		AdminViewed:     pitchRequest.AdminViewed,
		CreatedAt:       pitchRequest.CreatedAt.Unix(),
		UpdatedAt:       pitchRequest.UpdatedAt.Unix(),
	}

	// I am susppoed to return the resource Location in the header of the response, 
	// I haven't learned that and I will adhere in all REST projects from now on.
	ctx.JSON(http.StatusCreated, resp)
}

type ViewPitchReq struct {
	SalesRepID int64 `form:"sales_rep_id" binding:"required"`
	PageID     int32 `form:"page_id" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// salesViewPitchRequest api endpoint for getting pitchrequests by the sales rep
// on update to this function, the user who calls the function must be the owner of the sales_rep_id to be authorized
func (s *Server) salesViewPitchRequests(ctx *gin.Context) {
	var req ViewPitchReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !authAccess(ctx, utils.SalesRole) {
		return
	}
	args := db.ViewPitchRequestsParams{
		SalesRepID: req.SalesRepID,
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
	}

	requests, err := s.Store.ViewPitchRequests(ctx, args)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, requests)
}

// DeletePitchReq holds fields required to delete pitch request
type DeletePitchReq struct {
	ID         int64 `uri:"pitch_id" binding:"required"`
	SalesRepID int64 `uri:"sales_rep_id" binding:"required"`
	Username string `uri:"sales_username" binding:"required"`
}


// salesDeletePitchReqHandler api handler for deleting a pitch request
// username must match auth payload username.... (in any update, that bug should be fixed,
// when a user is updated and the username changes, the session should be refreshed if the user is logged in
func (s *Server) salesDeletePitchReqHandler(ctx *gin.Context) {
	var req DeletePitchReq

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authPayloadKey).(*token.Payload)
	if payload.Role != utils.SalesRole || payload.Username != req.Username {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}
	// Check if the pitch request exists
	args := db.PitchRequestExistParams{
		ID:         req.ID,
		SalesRepID: req.SalesRepID,
	}
	exists, err := s.Store.PitchRequestExist(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if !exists {
		ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("pitch request doesn't exist")))
		return
	}
	err = s.Store.DeletePitchRequest(ctx, req.ID)
	if err != nil {
		if pqErrHandler(ctx, "pitch request", err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{
		"message": "successful",
	})
}


type SalesDealsReq struct {
	SalesRepName string `form:"sales_rep" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required"`
	PageID int32 `form:"page_id" binding:"required"`
}

// getSalesDeals returns deals associated with a sales rep pitch request
func (s *Server) getSalesDeals(ctx *gin.Context) {
	var req SalesDealsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !authAccess(ctx, utils.SalesRole) {
		return
	}
	args := db.GetDealsBySalesRepParams{
		SalesRepName: req.SalesRepName,
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	deals, err := s.Store.GetDealsBySalesRep(ctx, args)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, deals)
}