package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

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
	CustomerRequest string   `json:"customer_request" binding:"required"`
	RequestDeadline UnixTime `json:"request_deadline" binding:"required"`
}

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
		CustomerRequest: req.CustomerRequest,
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
		CustomerRequest: pitchRequest.CustomerRequest,
		RequestDeadline: pitchRequest.RequestDeadline,
		AdminViewed:     pitchRequest.AdminViewed,
		CreatedAt:       pitchRequest.CreatedAt,
		UpdatedAt:       pitchRequest.UpdatedAt.Time,
	}
	ctx.JSON(http.StatusOK, resp)
}

type SalesUpdateUserReq struct {
	ID       int64  `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required,alphanum"`
}

// LoginReq holds fields required to access user details
func (s *Server) salesUpdateuserHandler(ctx *gin.Context) {
	var req SalesUpdateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !authAccess(ctx, utils.SalesRole) {
		return
	}
	//get user for update
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
	args := db.UpdateUserParams{
		ID:        usr.ID,
		Username:  req.Username,
		UpdatedAt: updateTime,
	}
	// get
	newUsr, err := s.Store.UpdateUser(ctx, args)
	if err != nil {
		if pqErrHandler(ctx, "sales-rep", err) {
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
		PasswordChanged: newUsr.PasswordChanged,
		UpdatedAt:       newUsr.UpdatedAt.Time,
		CreatedAt:       newUsr.CreatedAt,
	}
	ctx.JSON(http.StatusOK, resp)
}

type ViewPitchReq struct {
	SalesRepID int64 `form:"sales_rep_id" binding:"required"`
	PageID     int32 `form:"page_id" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=5,max=10"`
}

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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, requests)
}

type DeletePitchReq struct {
	ID         int64 `uri:"pitch_id" binding:"required"`
	SalesRepID int64 `uri:"sales_rep_id" binding:"required"`
	Username string `uri:"sales_username" binding:"required"`
}

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

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successful",
	})
}


type SalesDealsReq struct {
	SalesRepName string `json:"sales_rep" binding:"required"`
	PageSize int32 `json:"page_size" binding:"required"`
	PageID int32 `json:"page_id" binding:"required"`
}

// getSalesDeals returns deals associated with a sales rep pitch request
func (s *Server) getSalesDeals(ctx *gin.Context) {
	var req SalesDealsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
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