package api

import (
	"fmt"
	"log/slog"

	"github.com/ekefan/pre-sales-deal-tracker/backend/middleware"
	db "github.com/ekefan/pre-sales-deal-tracker/backend/db/sqlc"
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
	tokenGen, err := token.NewPasetoGenerator(server.config.SymmetricKey)
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
	authGrp.POST("/users", server.createUser)
	server.router = router
	slog.Info("Router is setup and ready to run")

}

func (server *Server) StartServer() error {
	initUser(server.store)
	return server.router.Run(server.config.ServerAddres)
}
