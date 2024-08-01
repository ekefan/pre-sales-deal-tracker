package server

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// CreateUsrReq holds fields that must be provided by client to create user
type CreateUsrReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Role     string `json:"Role" binding:"required"`
	FullName string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// CreateUsrRep holds fields that must be provided to client after creating a user
type CreateUsrResp struct {
	Role      string    `json:"Role"`
	ID        int64     `json:"user_id"`
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

	//hashpassword

	args := db.CreateNewUserParams{
		Username: req.Username,
		Role:     req.Role,
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := s.Store.CreateNewUser(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := CreateUsrResp{
		Role:      user.Role,
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, resp)
}

// AdminUpdateUsrReq holds the field - ID to unmarshall json requests
// ID here is the id of the user to be updated
// They are all required, however if no new values are passed... the current
// the current user fields will be used
type AdminUpdateUsrReq struct {
	ID       int64  `json:"user_id" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required,alphanum"`
}

// AdminUpdateUsrResp holds the fields for responding accurately to updating user end-point
type AdminUpdateUsrResp struct {
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// adminUpdateUserHandler http handler for the api end point for updating a user
func (s *Server) adminUpdateUserHandler(ctx *gin.Context) {
	var req AdminUpdateUsrReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//get user for update in the kitch
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
	args := db.AdminUpdateUserParams{
		ID:        usr.ID,
		FullName:  req.Fullname,
		Email:     req.Email,
		Password:  req.Password,
		Username:  req.Username,
		UpdatedAt: updateTime,
	}
	// get
	newUsr, err := s.Store.AdminUpdateUser(ctx, args)
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

// AdminDeleteUserReq holds field user id that is to be deleted
type AdminDeleteUserReq struct {
	ID int64 `uri:"id" binding:"required"`
}

// adminDeleteUserhandler http handler for the api end point for Deleting a user
func (s *Server) adminDeleteUserHandler(ctx *gin.Context) {
	var req AdminDeleteUserReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.Store.AdminDeleteUser(ctx, req.ID)
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
		PitchID:             req.PitchID,
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
		PitchID:             deal.PitchID,
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
	ID                  int64     `json:"id" binding:"required"`
	ServiceToRender     string    `json:"service_to_render" binding:"required"`
	Status              string    `json:"status" binding:"required"`
	StatusTag           string    `json:"statusTag" binding:"required"`
	CurrentPitchRequest string    `json:"current_pitch_request" binding:"required"`
	UpdatedAt           time.Time `json:"updated_at" binding:"required"`
	ClosedAt            time.Time `json:"closed_at" binding:"required"`
}

// adminUpdateDealsHandler http handler for the api end point for updating a deal
func (s *Server) adminUpdateDealHandler(ctx *gin.Context) {
	var req UpdateDealReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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
		Time:  req.UpdatedAt,
		Valid: true,
	}
	closedAt := sql.NullTime{
		Time:  req.UpdatedAt,
		Valid: !req.UpdatedAt.IsZero(),
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
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := CreateDealResp{
		ID:                  updatedDeal.ID,
		PitchID:             updatedDeal.PitchID,
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
	ID int64 `uri:"id" binding:"required"`
}

// adminDeleteDealHandler http handler for the api end point for Deleting a deal
func (s *Server) adminDeleteDealHandler(ctx *gin.Context) {
	var req DeleteDealReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.Store.AdminDeleteDeal(ctx, req.ID)
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

type ListUsersReq struct {
	PageID   int32 `json:"page_id" binding:"required"`
	PageSize int32 `json:"page_size" binding:"required"`
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
