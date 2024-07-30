package main

import (
	"database/sql"
	"log"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/ekefan/deal-tracker/internal/server"
	_ "github.com/lib/pq"
)

const (
	dbDriver string = "postgres"
	dbSource string = "postgresql://root:vasDealTracker@localhost:5432/dealTrackerDB?sslmode=disable"
)

// main entry point of the application
func main(){
	//connect to database
	dbConn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("couldn't connect to database", err)
	}
	//Create new store instance
	store := db.NewStore(dbConn)
	//Create new server instance
	server := server.NewServer(store)
	//Setup router
	server.SetupRouter()
	//Start Server
	if err := server.StartServer("0.0.0.0:8080"); err != nil {
		log.Fatal("couldn't start server", err)
	}
}
