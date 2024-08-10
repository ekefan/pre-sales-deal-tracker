package server

import (
	"fmt"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/ekefan/deal-tracker/internal/token"
	"github.com/ekefan/deal-tracker/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server contains fields required by the server instance
type Server struct {
	Router *gin.Engine
	Store  db.Store
	EnvVar Config
	TokenMaker token.TokenMaker
}

// NewServer create a server instance, having a router that connect api endpoints
func NewServer(store db.Store, config Config) (*Server, error) {
	tokenMaker, err := token.NewPaseto(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("couln't create token: %w", err)
	}
	return &Server{
		Store: store,
		EnvVar: config,
		TokenMaker: tokenMaker,
	}, nil
}

// SetupRouter ini
func (s *Server) SetupRouter() {
	router := gin.Default()

	// register validation to gin context
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("valid-role", utils.RoleValidator)
	}

	// ADMIN
	router.POST("/users", s.adminCreateUserHandler)
	router.POST("/users/login", s.userLogin) //added token

	authRoute := router.Group("/auth").Use(authMiddleware(s.TokenMaker))
	authRoute.PUT("/users/update/", s.adminUpdateUserHandler) //added token authorization
	router.DELETE("/users/delete/:id/admin_role", s.adminDeleteUserHandler) //added token authorization
	router.POST("/admin/deals", s.adminCreateDealHandler) 
	router.PUT("admin/deals/update", s.adminUpdateDealHandler)
	router.DELETE("/admin/deals/delete/:deal_id/:admin_role", s.adminDeleteDealHandler)
	router.GET("/users", s.listUsersHandler)


	// ADMINSALES
	router.PUT("pitchrequest/update", s.updatePitchReqHandler)
	router.GET("/deals", s.getDealsHandler)
	router.GET("/deals/vas", s.getOngoingDeals)
	router.GET("/deals/filtered", s.getFilteredDeals)
	router.GET("deals/filtered/count", s.getCountFilteredDeals)

	//SALES-REP
	router.POST("/sales/pitchReq", s.salesCreatePitchReqHandler)
	router.PUT("/sales/update/:username", s.salesUpdateuserHandler)
	router.GET("/pitchrequest/", s.salesViewPitchRequests)
	router.DELETE("/sales/pitchReq/delete/:sales_rep_id/:pitch_id", s.salesDeletePitchReqHandler)
	router.GET("sales/deals", s.getSalesDeals)
	// router.GET("sales/count/deals", s.getSalesDealsCount)

	//General
	router.PUT("/users/password", s.updatePassWordLoggedIn)
	router.PUT("/users/forgotpassword", s.forgotPassword)

	s.Router = router
}

// hostAddress string
func (s *Server) StartServer(hostAddress string) error {
	err := s.Router.Run(hostAddress)
	if err != nil {
		return err
	}
	return nil
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
