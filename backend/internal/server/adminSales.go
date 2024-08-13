package server

////////// handlers /////////

// userLogin
//
import (
	"fmt"
	"net/http"
	"time"

	// db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/ekefan/deal-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

// LoginReq holds fields required to access user details
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	ID        int64     `json:"user_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	FullName  string    `json:"fullname"`
	Email     string    `json:"email"`
	UpdatedAt int64 `json:"updated_at"`
	CreatedAt int64 `json:"created_at"`
}

type UserResp struct {
	AccessToken string    `json:"access_token"`
	LoginResp   LoginResp `json:"user"`
}

// userLogin handler for loging users to application
func (s *Server) userLogin(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check if the user has been created
	user, err := s.Store.GetUser(ctx, req.Username)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//this is frontend logic
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("password invalid: %v", err)))
		return
	}

	// create accessToken
	accessToken, err := s.TokenMaker.CreateToken(user.Username, user.Role, s.EnvVar.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//Create token or cookies and add to resp
	usr := LoginResp{
		//accessToken:
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		FullName:  user.FullName,
		Email:     user.Email,
		UpdatedAt: user.UpdatedAt.Unix(),
		CreatedAt: user.CreatedAt.Unix(),
	}
	resp := UserResp{
		AccessToken: accessToken,
		LoginResp:   usr,
	}
	ctx.JSON(http.StatusOK, resp)

}

type UpdatePitchReq struct {
	ID              int64    `json:"pitch_request_id" binding:"required"`
	Status          string   `json:"status" binding:"required"`
	PitchTag        string   `json:"pitch_tag" binding:"required"`
	CustomerRequests []string   `json:"customer_requests" binding:"required"`
	AdminViewed     bool     `json:"admin_viewed"`
	RequestDealine  UnixTime `json:"request_deadline"`
}

type PitchResp struct {
	ID              int64     `json:"pitch_request_id"`
	SalesRepID      int64     `json:"sales_rep_id"`
	SalesRepName    string    `json:"sales_rep_name"`
	Status          string    `json:"status"`
	CustomerName    string    `json:"customer_name"`
	PitchTag        string    `json:"pitch_tag"`
	CustomerRequests []string    `json:"customer_requests"`
	RequestDeadline int64 `json:"request_deadline"`
	AdminViewed     bool      `json:"admin_viewed"`
	CreatedAt       int64 `json:"created_at"`
	UpdatedAt       int64 `json:"updated_at"`
}

func (s *Server) updatePitchReqHandler(ctx *gin.Context) {
	var req UpdatePitchReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// getAuthPayload and require role == admin or sales
	if !multipleAuthAccess(ctx, []string{utils.AdminRole, utils.SalesRole}) {
		return
	}

	pitchReq, err := s.Store.GetPitchRequestForUpdate(ctx, req.ID)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var deadline time.Time
	if req.RequestDealine.Valid {
		deadline = req.RequestDealine.Time
	} else {
		deadline = pitchReq.RequestDeadline
	}
	args := db.UpdatePitchRequestParams{
		ID:              pitchReq.ID,
		Status:          req.Status,
		PitchTag:        req.PitchTag,
		CustomerRequest: req.CustomerRequests,
		AdminViewed:     req.AdminViewed,
		UpdatedAt:       time.Now(),
		RequestDeadline: deadline,
	}
	updatedPitchReq, err := s.Store.UpdatePitchRequest(ctx, args)
	if err != nil {
		if pqErrHandler(ctx, "pitch_requests", err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := PitchResp{
		ID:              updatedPitchReq.ID,
		SalesRepID:      updatedPitchReq.SalesRepID,
		SalesRepName:    updatedPitchReq.SalesRepName,
		Status:          updatedPitchReq.Status,
		CustomerName:    updatedPitchReq.CustomerName,
		PitchTag:        updatedPitchReq.PitchTag,
		CustomerRequests: updatedPitchReq.CustomerRequest,
		RequestDeadline: updatedPitchReq.RequestDeadline.Unix(),
		AdminViewed:     updatedPitchReq.AdminViewed,
		CreatedAt:       updatedPitchReq.CreatedAt.Unix(),
		UpdatedAt:       updatedPitchReq.UpdatedAt.Unix(),
	}
	ctx.JSON(http.StatusOK, resp)
}
