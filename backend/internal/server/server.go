package server

import (
	"fmt"
	"time"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/ekefan/deal-tracker/internal/token"
	"github.com/ekefan/deal-tracker/internal/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server contains fields required by a server instance
type Server struct {
	Router     *gin.Engine      // Router an instance of gin.Engine
	Store      db.Store         // Store the database interface for interacting with the db
	EnvVar     Config           // EnvVar holds the environment variables loaded into the server instance
	TokenMaker token.TokenMaker // interface for creating and managing tokens
}

// NewServer create a server instance, having a router that connect api endpoints
func NewServer(store db.Store, config Config) (*Server, error) {
	tokenMaker, err := token.NewPaseto(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("couln't create token: %w", err)
	}
	return &Server{
		Store:      store,
		EnvVar:     config,
		TokenMaker: tokenMaker,
	}, nil
}

// SetupRouter sets up a router, register cors, and routes for api endpoints
func (s *Server) SetupRouter() {
	router := gin.Default()
	// register validation to gin context
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("valid-role", utils.RoleValidator)
	}
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(corsConfig))
	authRoute := router.Group("/").Use(authMiddleware(s.TokenMaker))

	router.POST("/users", s.adminCreateUserHandler)
	router.POST("/users/login", s.userLogin)
	authRoute.POST("/sales/pitchReq", s.salesCreatePitchReqHandler)
	authRoute.POST("/admin/deals", s.adminCreateDealHandler)
	authRoute.GET("/deals/vas", s.getOngoingDeals)
	authRoute.GET("/deals/filtered", s.getFilteredDeals)
	authRoute.GET("/admin/pitchrequest", s.adminGetPitchRequests)
	authRoute.GET("/sales/pitchrequest", s.salesViewPitchRequests)
	authRoute.GET("/admin/getdeal", s.getDealsById)
	authRoute.GET("/sales/deals", s.getSalesDeals)
	// FIXME: not complain with REST-API standards => GET /users
	authRoute.GET("/list-users", s.listUsersHandler)
	authRoute.PUT("/users/password-reset", s.resetPassword)
	authRoute.PUT("/admin/deals/update", s.adminUpdateDealHandler)
	authRoute.PUT("/users/update", s.adminUpdateUserHandler)
	authRoute.PUT("/users/password", s.updatePassWordLoggedIn)
	authRoute.PUT("/pitchrequest/update", s.updatePitchReqHandler)
	authRoute.DELETE("/users/delete/:id", s.adminDeleteUserHandler)
	authRoute.DELETE("/admin/deals/delete/:deal_id", s.adminDeleteDealHandler)
	authRoute.DELETE("/sales/pitchReq/delete/:sales_username/:sales_rep_id/:pitch_id", s.salesDeletePitchReqHandler)

	s.Router = router
}

// StartServer starts the app server, takes the hostAddress
// listens and serves on that address
func (s *Server) StartServer(hostAddress string) error {
	err := s.Router.Run(hostAddress)
	if err != nil {
		return err
	}
	return nil
}

// errorResponse a custom error reponse handler for reusability
// takes the  error (err) returns a gin.H{} struct with an error
// field equal to err.Error()
func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
