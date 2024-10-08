package api

import (
	"fmt"
	"log/slog"
	"time"

	db "github.com/ekefan/pre-sales-deal-tracker/backend/internal/db/sqlc"
	"github.com/ekefan/pre-sales-deal-tracker/backend/middleware"
	"github.com/ekefan/pre-sales-deal-tracker/backend/token"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server servers http requests for the deal tracker service
type Server struct {
	store          db.Store
	config         *Config
	tokenGenerator token.TokenGenerator
	router         *gin.Engine
}

// NewServer creates a token generator and sets up a router with store and config
func NewServer(store db.Store, config *Config) (*Server, error) {
	server := &Server{
		store:  store,
		config: config,
	}

	// set token generator
	// [Q]: let's have a discussion around this Paseto thing since I've never used it.
	tokenGen, err := token.NewJwtGenerator(server.config.SymmetricKey)
	if err != nil {
		fmt.Println(len(server.config.SymmetricKey))
		return nil, fmt.Errorf("cannot generate tokens, %v", err)
	}

	server.tokenGenerator = tokenGen
	server.setupRouter()
	return server, nil
}

// setupRouter uses a default gin engine
// defines api endpoints, middlewares and handler functions
func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                             // Allow all origins or specify origins
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},        // Specify allowed methods
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // Specify allowed headers
		ExposeHeaders:    []string{"Content-Length"},                // Optional: specify headers that can be exposed
		AllowCredentials: true,                                      // Optional: allow credentials
		MaxAge:           12 * time.Hour,                            // Optional: set max age for preflight requests
	}))
	router.POST("/auth/login", server.authLogin)
	// for authentication middlewares
	authGrp := router.Group("/")

	// for authorization middlwares
	adminGrp := router.Group("/")
	salesGrp := router.Group("/")

	authGrp.Use(middleware.UserAuthentication(server.tokenGenerator))

	adminGrp.Use(middleware.UserAuthentication(server.tokenGenerator))
	adminGrp.Use(middleware.AdminAccessAuthorization())

	salesGrp.Use(middleware.UserAuthentication(server.tokenGenerator))
	salesGrp.Use(middleware.SalesAccessAuthorization())

	authGrp.GET("/deals", server.retrieveDeals)
	authGrp.PUT("/pitch_requests/:pitch_id", server.updatePitchReq)

	adminGrp.POST("/users", server.createUser)
	adminGrp.GET("/users", server.retrieveUsers)
	authGrp.GET("/users/:user_id", server.getUsersByID)
	adminGrp.PUT("/users/:user_id", server.updateUser)
	adminGrp.DELETE("/users/:user_id", server.deleteUser)
	adminGrp.PATCH("/users/:user_id/password", server.updateUserPassword)
	adminGrp.POST("/deals", server.createDeal)
	adminGrp.PUT("/deals/:deal_id", server.updateDeal)
	adminGrp.DELETE("/deals/:deal_id", server.deleteDeal)

	salesGrp.POST("/pitch_requests", server.createPitchRequest)
	salesGrp.GET("/pitch_requests", server.retrievePitchRequests)
	salesGrp.DELETE("/pitch_requests/:pitch_id", server.deletePitchRequest)

	server.router = router
	slog.Info("Router is setup and ready to run")
}

func (server *Server) StartServer() error {
	initUser(server.store)
	return server.router.Run(server.config.ServerAddres)
}
