package server

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// PitchReq holds fields for creating a Pitch Request
type PitchReq struct {
	SalesRepID      int64     `json:"sales_rep_id" binding:"required,min=6"`
	SalesRepName    string    `json:"sales_rep_name" binding:"required"`
	Status          string    `json:"status" binding:"required"`
	CustomerName    string    `json:"customer_name" binding:"required"`
	PitchTag        string    `json:"pitch_tag" binding:"required"`
	CustomerRequest string    `json:"customer_request" binding:"required"`
	RequestDeadline time.Time `json:"request_deadline" binding:"required"`
}

func (s *Server) salesCreatePitchReqHandler(ctx *gin.Context) {
	var req PitchReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// instead of receiving salesRepID from json validate it through payload
	args := db.CreatePitchRequestParams{
		SalesRepID:      req.SalesRepID,
		SalesRepName:    req.SalesRepName,
		Status:          req.Status,
		CustomerName:    req.CustomerName,
		PitchTag:        req.PitchTag,
		CustomerRequest: req.CustomerRequest,
		RequestDeadline: req.RequestDeadline,
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
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required,alphanum"`
}

// LoginReq holds fields required to access user details
func (s *Server) salesUpdateuserHandler(ctx *gin.Context) {
	var req SalesUpdateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//get user for update,  get ID from authorization payload
	usr, err := s.Store.GetUserForUpdate(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	// Set update time to time now....
	updateTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	// Hash password
	args := db.UpdateUserParams{
		ID:        usr.ID,
		Password:  req.Password,
		Username:  req.Username,
		UpdatedAt: updateTime,
	}
	// get
	newUsr, err := s.Store.UpdateUser(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	resp := AdminUpdateUsrResp{
		UserID:    newUsr.ID,
		Username:  newUsr.Username,
		Role:      newUsr.Role,
		Fullname:  newUsr.FullName,
		Email:     newUsr.Email,
		UpdatedAt: newUsr.UpdatedAt.Time,
		CreatedAt: newUsr.CreatedAt,
	}
	ctx.JSON(http.StatusOK, resp)
}

type ViewPitchReq struct {
	ID       int64 `uri:"pitch_id" binding:"required"`
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) salesViewPitchRequests(ctx *gin.Context) {
	var req ViewPitchReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args := db.ViewPitchRequestsParams{
		ID:     req.ID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	requests, err := s.Store.ViewPitchRequests(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, requests)
}


type DeletePitchReq struct {
	ID int64 `uri:"id" binding:"required"`
}
func (s *Server) salesDeletePitchReqHandler(ctx *gin.Context) {
	var req DeletePitchReq

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := s.Store.DeletePitchRequest(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successful",
	})
}