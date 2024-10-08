package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	db "github.com/ekefan/pre-sales-deal-tracker/backend/internal/db/sqlc"
	"github.com/ekefan/pre-sales-deal-tracker/backend/middleware"
	"github.com/ekefan/pre-sales-deal-tracker/backend/token"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreatePitchReq holds the fields needed to create a pitch request
type CreatePitchReq struct {
	UserId           int64     `json:"user_id"`
	CustomerName     string    `json:"customer_name"`
	CustomerRequests []string  `json:"customer_request"`
	AdminTask        string    `json:"admin_task"`
	AdminDeadline    time.Time `json:"admin_deadline"`
}

// createPitchRequest
func (server *Server) createPitchRequest(ctx *gin.Context) {
	var req CreatePitchReq
	if err := bindClientRequest(ctx, &req, jsonSource); err != nil {
		handleClientReqError(ctx, err)
		return
	}
	pitchReqCreated, err := server.store.CreatePitchRequest(ctx, db.CreatePitchRequestParams{
		UserID:          req.UserId,
		CustomerName:    req.CustomerName,
		CustomerRequest: req.CustomerRequests,
		AdminTask:       req.AdminTask,
		AdminDeadline:   pgtype.Timestamp{Time: req.AdminDeadline, Valid: true},
	})
	if err != nil {
		details := fmt.Sprintf("could not create pitch request because no user with id: %v, exists", req.UserId)
		if handleDbError(ctx, err, details) {
			return
		}
		handleServerError(ctx, err)
		return
	}
	if pitchReqCreated < 1 {
		handleServerError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, successMessage())
}

// retrievePitchRequests end point handler for get /pitch_requests,
// gets all pitch requests created by the authorized user
func (server *Server) retrievePitchRequests(ctx *gin.Context) {
	var req GetPaginatedReq
	if err := bindClientRequest(ctx, &req, querySource); err != nil {
		handleClientReqError(ctx, err)
		return
	}
	payload := ctx.MustGet(middleware.AuthPayloadKey).(*token.Payload)
	result, err := server.store.GetPitchRequestsPaginated(ctx, db.GetPitchRequestsPaginatedParams{
		UserID: payload.UserID,
		Limit:  req.PageSize,
		Offset: req.PageSize * (req.PageID - 1),
	})
	if err != nil {
		slog.Error(err.Error())
		handleServerError(ctx, err)
		return
	}
	PitchRequestData := []PitchRequest{}
	totalPitchRequests := result.TotalPitchRequests
	if len(result.PitchRequests) > 0 {
		err := json.Unmarshal(result.PitchRequests, &PitchRequestData)
		if err != nil {
			slog.Error(err.Error())
			handleServerError(ctx, err)
			return
		}
	}
	resp := struct {
		Pagination `json:"pagination"`
		Data       []PitchRequest `json:"data"`
	}{
		Pagination: generatePagination(int32(totalPitchRequests), req.PageID, req.PageSize),
		Data:       PitchRequestData,
	}
	ctx.JSON(http.StatusOK, resp)
}

// PitchIDFromURI holds the pitch_id for needed to access a pitch_request
type PitchIDFromURI struct {
	PitchID int64 `uri:"pitch_id" binding:"required,min=1,numeric"`
}

// UpdatePitchReq holds the fields needed to update a pitch request
type UpdatePitchReq struct {
	AdminViewed     bool     `json:"admin_viewed"`
	CustomerRequest []string `json:"customer_request"`
}

// updatePitchReq handler for put /pitch_requests/:pitch_id, updates a pitch request
func (server *Server) updatePitchReq(ctx *gin.Context) {
	var (
		reqUri  PitchIDFromURI
		reqBody UpdatePitchReq
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

	numPitchReqUpdated, err := server.store.UpdatePitchRequest(ctx, db.UpdatePitchRequestParams{
		AdminViewed:     reqBody.AdminViewed,
		CustomerRequest: reqBody.CustomerRequest,
		PitchID:         reqUri.PitchID,
	})
	if err != nil {
		details := fmt.Sprintf("can not update pitch request with id: %v", reqUri.PitchID)
		if handleDbError(ctx, err, details) {
			return
		}
		fmt.Println("here", err.Error())
		handleServerError(ctx, err)
		return
	}
	if numPitchReqUpdated != 1 {
		details := fmt.Sprintf("no pitch request with id: %v, exists", reqUri.PitchID)
		handleDbError(ctx, errNotFound, details)
		return
	}
	ctx.JSON(http.StatusOK, successMessage())
}

// deletePitchRequest handler for delete /pitch_request/:pitch_id, deletes a pitch request
func (server *Server) deletePitchRequest(ctx *gin.Context) {
	var req PitchIDFromURI
	if err := bindClientRequest(ctx, &req, uriSource); err != nil {
		handleClientReqError(ctx, err)
		return
	}
	numPitchReqDeleted, err := server.store.DeletePitchRequest(ctx, req.PitchID)
	if err != nil {
		slog.Error(err.Error())
		handleServerError(ctx, err)
		return
	}
	if numPitchReqDeleted < 1 {
		err := errNotFound
		detail := err.Error()
		handleDbError(ctx, err, detail)
		return
	}
	ctx.Status(http.StatusNoContent)
}
