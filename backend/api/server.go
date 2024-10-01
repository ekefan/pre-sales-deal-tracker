package api

import (
	"fmt"
	"log/slog"

	db "github.com/ekefan/pre-sales-deal-tracker/backend/db/sqlc"
	"github.com/ekefan/pre-sales-deal-tracker/backend/token"
	"github.com/gin-gonic/gin"
)

// Server servers http requests for the deal tracker service
type Server struct {
	store          db.Store
	tokenGenerator token.TokenGenerator
	router         *gin.Engine
}

// NewServer creates a new http server and sets up  routing
func NewServer(store db.Store) (*Server, error) {
	tokenGen, err := token.NewPasetoGenerator("01234567890123456789012345678909")
	if err != nil {
		return nil, fmt.Errorf("cannot generate tokens: %v", err)
	}
	server := &Server{
		store: store,
		tokenGenerator: tokenGen,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/auth/login", server.authLogin)
	server.router = router
	slog.Info("Router is setup and ready to run")

}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}
