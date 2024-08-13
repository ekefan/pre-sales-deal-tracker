package main

import (
	"database/sql"
	"log"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/ekefan/deal-tracker/internal/server"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// main entry point of the application
func main() {
	//connect to database
	gin.SetMode(gin.TestMode)
	config, err := server.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}
	dbConn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("couldn't connect to database", err)
	}
	//Create new store instance
	store := db.NewStore(dbConn)
	//Create new server instance
	server, err := server.NewServer(store, config)
	if err != nil {
		log.Fatal("could't spin up server: %w", err)
	}
	//Setup router
	server.SetupRouter()
	//Start Server
	if err := server.StartServer("0.0.0.0:8080"); err != nil {
		log.Fatal("couldn't start server", err)
	}
}
