package server

// handlers
// adminCreateDealHandler
// adminUpdateDealHandler
// adminDeleteDealHandler
// listUsersHandler
import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/ekefan/deal-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

// CreateDealReq holds fields needed to create a new deal
type CreateDealReq struct {
	PitchID int64 `json:"pitch_id" binding:"required,gt=0"`
}

type CreateDealResp struct {
	ID                  int64    `json:"id"`
	PitchID             int64    `json:"pitch_id"`
	SalesRepName        string   `json:"sales_rep_name"`
	CustomerName        string   `json:"customer_name"`
	ServiceToRender     []string `json:"service_to_render"`
	Status              string   `json:"status"`
	StatusTag           string   `json:"status_tag"`
	CurrentPitchRequest string   `json:"current_pitch_request"`
	NetTotalCost        string   `json:"net_total_cost"`
	Profit              string   `json:"profit"`
	CreatedAt           int64    `json:"created_at"`
	UpdatedAt           int64    `json:"updated_at"`
	ClosedAt            int64    `json:"closed_at"`
	Awarded             bool     `json:"awarded"`
}

// adminCreateDealHandler http handler for the api end point for creating a new deal
func (s *Server) adminCreateDealHandler(ctx *gin.Context) {
	// validate json-req and unmarshall it to the req
	var req CreateDealReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// find a way to perform authorization
	if !authAccess(ctx, utils.AdminRole) {
		return
	}

	// make call to database
	success := s.Store.CreateDealTxn(ctx, db.CreateDealTxnArgs{
		PitchID: req.PitchID,
	})

	if success != nil {
		if pqErrHandler(ctx, "deal", success) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(success))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"mesaage": "succesful",
	})
}

// UpdateDealReq holds fields used to update a deal
type UpdateDealReq struct {
	ID                  int64    `json:"id" binding:"required,numeric"`
	ServicesToRender    []string `json:"services_to_render" binding:"required"`
	Status              string   `json:"status" binding:"required"`
	StatusTag           string   `json:"department" binding:"required"`
	CurrentPitchRequest string   `json:"current_pitch_request" binding:"required"`
	NetTotalCost        int64    `json:"net_total_cost" binding:"numeric,gte=0"`
	Profit              int64    `json:"profit" binding:"numeric,gte=0"`
	Awarded             bool     `json:"awarded" binding:"boolean"`
}

// adminUpdateDealsHandler http handler for the api end point for updating a deal
func (s *Server) adminUpdateDealHandler(ctx *gin.Context) {
	var req UpdateDealReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// verify that the resouces is accessed by the admin only, check for role in
	// authorization payload
	if !authAccess(ctx, utils.AdminRole) {
		return
	}

	deal, err := s.Store.AdminGetDealForUpdate(ctx, req.ID)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) || pqErrHandler(ctx, "deals", err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var closedAt time.Time
	if req.Awarded {
		closedAt = time.Now()
	} else {
		closedAt = deal.ClosedAt
	}
	netTotal := strconv.Itoa(int(req.NetTotalCost))
	profit := strconv.Itoa(int(req.Profit))
	updatedDeal, err := s.Store.AdminUpdateDeal(ctx, db.AdminUpdateDealParams{
		ID:                  deal.ID,
		ServiceToRender:     req.ServicesToRender,
		Status:              req.Status,
		StatusTag:           req.StatusTag,
		CurrentPitchRequest: req.CurrentPitchRequest,
		UpdatedAt:           time.Now(),
		ClosedAt:            closedAt,
		Awarded:             req.Awarded,
		NetTotalCost:        netTotal,
		Profit:              profit,
	})
	if err != nil {
		if pqErrHandler(ctx, "deals", err) {
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
		NetTotalCost:        updatedDeal.NetTotalCost,
		Profit:              updatedDeal.Profit,
		CreatedAt:           updatedDeal.CreatedAt.Unix(),
		UpdatedAt:           updatedDeal.UpdatedAt.Unix(),
		ClosedAt:            updatedDeal.ClosedAt.Unix(),
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

	// get payload and check the for user.role
	if !authAccess(ctx, utils.AdminRole) {
		return
	}

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

// FIXME: this has nothing to do with the Deals resource. Should be put in the adminUser.go
// Also you should clarify the "admin" prefix meaning. Are those endpoints restricted? Or what's else?
// listUsershandler http handler for the api end point for getting list of users currently
func (s *Server) listUsersHandler(ctx *gin.Context) {
	var req ListUsersReq
	// FIXME: if I set page_id = 0 (or outside of the allowed boundaries), please adjust it to be a default value. page_id might also start from '0'. In case you go away from conventions, you need to be declarative and put it in the documentation.
	// "(Required) The page number or offset from which to start retrieving Users. Determines where the current page of results starts in the overall list." => doesn't state the range of allowed values
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !multipleAuthAccess(ctx, []string{utils.AdminRole, utils.ManagerRole}) {
		return
	}
	// FIXME: if you're accepting pagination info, you should return a paginated result, not only the collection of resources.
	args := db.AdminViewUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	users, err := s.Store.AdminViewUsers(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("i don't know")))
		return
	}
	ctx.JSON(http.StatusOK, users)
}

type AdminPitchReq struct {
	Admin_viewed bool `form:"admin_viewed" binding:"boolean"`
}

// adminGetPitchRequest api endpoint to get pitchrequests the admin has not viewed
func (s *Server) adminGetPitchRequests(ctx *gin.Context) {
	var req AdminPitchReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !authAccess(ctx, utils.AdminRole) {
		return
	}
	pitchRequests, err := s.Store.AdminGetPitchRequest(ctx, req.Admin_viewed)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("admin getting pitch requests I haven't fixed this yet %s", err)))
		return
	}
	ctx.JSON(http.StatusOK, pitchRequests)
}

// ===== TODO =====
///////////////////////////
// implement Handling user sessions
