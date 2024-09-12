package server

import (
	"fmt"
	"net/http"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/ekefan/deal-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type OngoingDealsReq struct {
	Status string `form:"status" binding:"required"`
}

func (s *Server) getOngoingDeals(ctx *gin.Context) {
	var req OngoingDealsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		fmt.Println("bad request", err, ctx.Request.Body)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !multipleAuthAccess(ctx, []string{utils.AdminRole, utils.ManagerRole, utils.SalesRole}) {
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
	CustomerName    *string  `form:"customer_name"`
	ServiceToRender []string `form:"service_to_render"`
	Status          *string  `form:"status"`
	MaxProfit       *string  `form:"max_profit"`
	MinProfit       *string  `form:"min_profit"`
	Awarded         *bool    `form:"awarded"`
	SalesRepName    *string  `form:"sales_rep_name"`
	PageSize        int32    `form:"page_size"`
	PageID          int32    `form:"page_id"`
}

func (s *Server) getFilteredDeals(ctx *gin.Context) {
	var req FilterDealReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !multipleAuthAccess(ctx, []string{utils.AdminRole, utils.ManagerRole, utils.SalesRole}) {
		return
	}
	args := db.FilterDealsParams{
		CustomerName:    req.CustomerName,
		ServiceToRender: req.ServiceToRender,
		Status:          req.Status,
		Profit:          req.MinProfit,
		Profit_2:        req.MaxProfit,
		Awarded:         req.Awarded,
		SalesRepName:    req.SalesRepName,
		Limit:           req.PageSize,
		Offset:          (req.PageID - 1) * req.PageSize,
	}

	deals, err := s.Store.FilterDeals(ctx, args)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		if pqErrHandler(ctx, "deals", err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, deals)

}


type GetDealReq struct {
	Deal_Id int64 `form:"deal_id" binding:"required"`
}

func (s *Server) getDealsById(ctx *gin.Context) {
	var req GetDealReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !authAccess(ctx, utils.AdminRole) {
		return
	}

	deal, err := s.Store.GetDealsById(ctx, req.Deal_Id)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		if pqErrHandler(ctx, "deals", err) {
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, deal)
}