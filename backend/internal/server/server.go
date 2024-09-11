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

// Server contains fields required by the server instance
type Server struct {
	Router     *gin.Engine
	Store      db.Store
	EnvVar     Config
	TokenMaker token.TokenMaker
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

// SetupRouter ini
func (s *Server) SetupRouter() {
	router := gin.Default()
	// register validation to gin context
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("valid-role", utils.RoleValidator)
	}
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow all origins
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(corsConfig))
	
	// ADMIN
	authRoute := router.Group("/").Use(authMiddleware(s.TokenMaker))
	router.POST("/users", s.adminCreateUserHandler)
	

	
	//Currently used routes in the application
	authRoute.GET("/deals/vas", s.getOngoingDeals)
	authRoute.GET("/deals/filtered", s.getFilteredDeals)
	authRoute.GET("/admin/pitchrequest", s.adminGetPitchRequests)
	authRoute.GET("/sales/pitchrequest/", s.salesViewPitchRequests)
	router.POST("/users/login", s.userLogin)
	authRoute.PUT("/admin/deals/update", s.adminUpdateDealHandler)
	authRoute.GET("/users", s.listUsersHandler)
	authRoute.PUT("/users/password-reset", s.resetPassword)


	authRoute.PUT("/users/update", s.adminUpdateUserHandler)                  
	authRoute.DELETE("/users/delete/:id", s.adminDeleteUserHandler)           
	authRoute.POST("/admin/deals", s.adminCreateDealHandler)                          
	authRoute.DELETE("/admin/deals/delete/:deal_id", s.adminDeleteDealHandler)
	                               
	

	// ADMINSALES
	authRoute.GET("/admin/getdeal", s.getDealsById)
	authRoute.PUT("pitchrequest/update", s.updatePitchReqHandler) 
	// this route can be replaced with filtered... or user makes call here first then filters
	authRoute.GET("/deals", s.getDealsHandler)                     
	authRoute.GET("/deals/filtered/count", s.getCountFilteredDeals)
	// handler exist in sales... updates users username only for sales
	authRoute.PUT("/sales/update/user", s.salesUpdateuserHandler)
	//SALES-REP
	authRoute.POST("/sales/pitchReq", s.salesCreatePitchReqHandler)

	//this should be the general update user without even for password
	//change password should be separate                                 (for sales only)
	                                                      
	authRoute.DELETE("/sales/pitchReq/delete/:sales_username/:sales_rep_id/:pitch_id", s.salesDeletePitchReqHandler) 
	authRoute.GET("/sales/deals", s.getSalesDeals)                                                                   
	// authRoute.GET("sales/count/deals", s.getSalesDealsCount)

	//Password Update
	authRoute.PUT("/users/password", s.updatePassWordLoggedIn) //added token authorization
	

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
