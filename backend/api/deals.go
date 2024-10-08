package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	db "github.com/ekefan/pre-sales-deal-tracker/backend/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateDealReq holds the pitch_id from which a deal is to be created
type CreateDealReq struct {
	PitchID int64 `json:"pitch_id" binding:"required,min=1,numeric"`
}

// createDeal end point handler for post /deals
func (server *Server) createDeal(ctx *gin.Context) {
	var req CreateDealReq
	if err := bindClientRequest(ctx, &req, jsonSource); err != nil {
		handleClientReqError(ctx, err)
		return
	}
	err := server.store.CreateDealFromPitchId(ctx, req.PitchID)
	if err != nil {
		details := fmt.Sprintf("couldn't create a deal for pitch request with id: %v", req.PitchID)
		if handleDbError(ctx, err, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, successMessage())
}

// retrieveDeals end point handler for get /deals
func (server *Server) retrieveDeals(ctx *gin.Context) {
	var req GetPaginatedReq
	if err := bindClientRequest(ctx, &req, querySource); err != nil {
		handleClientReqError(ctx, err)
		return
	}

	result, err := server.store.GetDealPaginated(ctx, db.GetDealPaginatedParams{
		Limit:  req.PageSize,
		Offset: req.PageSize * (req.PageID - 1),
	})
	if err != nil {
		slog.Error(err.Error())
		handleServerError(ctx, err)
		return
	}
	DealsData := []Deal{}
	totalDeals := result.TotalDeals
	if len(result.Deals) > 0 {
		dealErr := json.Unmarshal(result.Deals, &DealsData)
		if dealErr != nil {
			slog.Error(dealErr.Error())
			handleServerError(ctx, err)
			return
		}
	}
	resp := struct {
		Pagination `json:"pagination"`
		Data       []Deal `json:"data"`
	}{
		Data:       DealsData,
		Pagination: generatePagination(int32(totalDeals), req.PageID, req.PageSize),
	}
	ctx.JSON(http.StatusOK, resp)
}

// DealIDFromUri holds the deal_id needed to access a deal resource
type DealIDFromUri struct {
	DealID int64 `uri:"deal_id" binding:"required,min=1,numeric"`
}
// UpdateDealReq holds the fields needed to update a deal
type UpdateDealReq struct {
	ServicesToRender []string `json:"services_to_render"`
	Status           string   `json:"status"`
	Department       string   `json:"department"`
	NetTotalCost     float64  `json:"net_total_cost"`
	Profit           float64  `json:"profit"`
	Awarded          bool     `json:"awarded"`
}

// updateDealendpoint handler for put /deals/:deal_id, updates a deal
func (server *Server) updateDeal(ctx *gin.Context) {
	var (
		reqUri  DealIDFromUri
		reqBody UpdateDealReq
	)
	uriErr := bindClientRequest(ctx, &reqUri, uriSource)
	if uriErr != nil {
		handleClientReqError(ctx, uriErr)
		return
	}
	reqBodyErr := bindClientRequest(ctx, &reqBody, jsonSource)
	if reqBodyErr != nil {
		handleClientReqError(ctx, reqBodyErr)
		return
	}
	netTotal := pgtype.Numeric{}
	profit := pgtype.Numeric{}
	netTotalStr := strconv.FormatFloat(reqBody.NetTotalCost, 'f', -1, 64)
	profitStr := strconv.FormatFloat(reqBody.Profit, 'f', -1, 64)
	netTotal.Scan(netTotalStr)
	if err := netTotal.Scan(netTotalStr); err != nil {
		handleServerError(ctx, err)
		return
	}
	if err := profit.Scan(profitStr); err != nil {
		handleServerError(ctx, err)
		return
	}
	updated, err := server.store.UpdateDeals(ctx, db.UpdateDealsParams{
		ID: reqUri.DealID,
		ServicesToRender: reqBody.ServicesToRender,
		Status:           reqBody.Status,
		Department:       reqBody.Department,
		NetTotalCost:     netTotal,
		Profit:           profit,
		Awarded:          reqBody.Awarded,
	})
	if err != nil {
		details := fmt.Sprintf("can not update deal with id: %v", reqUri.DealID)
		if handleDbError(ctx, err, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	if updated != 1 {
		details := fmt.Sprintf("no deal with id: %v, exists", reqUri.DealID)
		handleDbError(ctx, errNotFound, details)
		return
	}
	ctx.JSON(http.StatusOK, successMessage())
}
// deleteDeal handler for delete /deals/:deal_id, deletes a deal
func (server *Server) deleteDeal(ctx *gin.Context) {
	var req DealIDFromUri
	if err := bindClientRequest(ctx, &req, uriSource); err != nil {
		handleClientReqError(ctx, err)
		return
	}

	numDealDeleted, err := server.store.DeleteDeals(ctx, req.DealID)
	if err != nil {
		slog.Error(err.Error())
		handleServerError(ctx, err)
		return
	}
	if numDealDeleted < 1 {
		err := errNotFound
		detail := err.Error()
		handleDbError(ctx, err, detail)
		return
	}
	ctx.Status(http.StatusNoContent)
}