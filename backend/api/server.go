package api

import (
	"fmt"
	"log/slog"

	db "github.com/ekefan/pre-sales-deal-tracker/backend/db/sqlc"
	"github.com/ekefan/pre-sales-deal-tracker/backend/middleware"
	"github.com/ekefan/pre-sales-deal-tracker/backend/token"
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
	router.POST("/auth/login", server.authLogin)

	authGrp := router.Group("/")
	authGrp.Use(middleware.UserAuthorization(server.tokenGenerator))

	adminGrp := router.Group("/")
	adminGrp.Use(middleware.UserAuthorization(server.tokenGenerator))
	adminGrp.Use(middleware.AdminAccessAuthorization())

	// salesGrp := router.Group("/")
	// salesGrp.Use(middleware.UserAuthorization(server.tokenGenerator))
	// salesGrp.Use(middleware.SalesAccessAuthorization())

	adminGrp.POST("/users", server.createUser)
	adminGrp.GET("/users", server.retrieveUsers)
	authGrp.GET("/users/:user_id", server.getUsersByID)
	adminGrp.PUT("/users/:user_id", server.updateUser)
	adminGrp.DELETE("/users/:user_id", server.deleteUsers)
	adminGrp.PATCH("/users/:user_id/password", server.updateUserPassword)
	server.router = router
	slog.Info("Router is setup and ready to run")
}

func (server *Server) StartServer() error {
	initUser(server.store)
	return server.router.Run(server.config.ServerAddres)
}
