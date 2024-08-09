package server

import (
	"net/http"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

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
	deals, err := s.Store.AdminViewDeals(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, deals)
}

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
	CustomerName    *string `json:"customer_name"`
	ServiceToRender *string `json:"service_to_render"`
	Status          *string `json:"status"`
	MaxProfit       *string `json:"max_profit"`
	MinProfit       *string `json:"min_profit"`
	Awarded         *bool   `json:"awarded"`
	SalesRepName    *string `json:"sales_rep_name"`
	PageSize        int32   `json:"page_size"`
	PageID          int32   `json:"page_id"`
}

func (s *Server) getFilteredDeals(ctx *gin.Context) {
	var req FilterDealReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.FilterDealsParams{
		Column1: req.CustomerName,    // column1 == customer_name
		Column2: req.ServiceToRender, // column2 == service_to_render
		Column3: req.Status,          // column3 == status
		Column4: req.MinProfit,       // column4 == min_profit
		Column5: req.MaxProfit,       // column5 == max_profit
		Column6: req.Awarded,         // column6 == awarded
		Column7: req.SalesRepName,    // column7 == sales_rep_name
		Limit:   req.PageSize,
		Offset:  (req.PageID - 1) * req.PageSize,
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

type CountFilterDealReq struct {
	CustomerName    *string `json:"customer_name"`
	ServiceToRender *string `json:"service_to_render"`
	Status          *string `json:"status"`
	MaxProfit       *string `json:"max_profit"`
	MinProfit       *string `json:"min_profit"`
	Awarded         *bool   `json:"awarded"`
	SalesRepName    *string `json:"sales_rep_name"`
}
func (s *Server) getCountFilteredDeals(ctx *gin.Context) {
	var req CountFilterDealReq 
	if err  := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CountFilteredDealsParams{
		Column1: req.CustomerName,    // column1 == customer_name
		Column2: req.ServiceToRender, // column2 == service_to_render
		Column3: req.Status,          // column3 == status
		Column4: req.MinProfit,       // column4 == min_profit
		Column5: req.MaxProfit,       // column5 == max_profit
		Column6: req.Awarded,         // column6 == awarded
		Column7: req.SalesRepName,    // column7 == sales_rep_name
	}
	count, err := s.Store.CountFilteredDeals(ctx, args)
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

	resp := struct{ Count int64 `json:"num_of_pages"`}{Count: count}
	ctx.JSON(http.StatusOK, resp)
}