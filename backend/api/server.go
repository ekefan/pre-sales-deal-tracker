package api

import (
	"log/slog"

	db "github.com/ekefan/pre-sales-deal-tracker/backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	
	server := &Server{
		store: store,
	}

	return server
}


func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/auth/login", server.authLogin)
	server.router = router
	slog.Info("Router is setup and ready to run")
	
}

func (server *Server) StartServer(address string ) error {
	server.setupRouter()
	return server.router.Run(address)
}

