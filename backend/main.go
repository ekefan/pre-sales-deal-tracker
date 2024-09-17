package main

import (
	"database/sql"
	"log"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/ekefan/deal-tracker/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// TODO: after you fixed the swagger.yml, you can try to visualize it in this website => https://editor.swagger.io/
// Be sure to import the file. You can also provide some default values. Given those defaults it's easier when you import the full collection in a tool like Postman.
// TODO: speaking about design, I think you should leave outside of the internal folder the HTTP Endpoints functions. Those endpoints are the things you may want to expose outside of your package. Then, all the technical implementation should be masked below the internal pkg.
// FIXME: another well-curated resource is this link of Zalando API documentation: https://opensource.zalando.com/restful-api-guidelines/#table-of-contents

// main entry point of the application
func main() {
	gin.SetMode(gin.TestMode)

	// load environment variables
	config, err := server.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}
	// connect to database
	dbConn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("couldn't connect to database", err)
	}

	// run migration if tables has not been created in the database
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		log.Fatal("couldn't migrate schema to database", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Fatal("couldn't migrate schema to database", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("couldn't migrate schema to database", err)
	}

	// Create new store instance
	store := db.NewStore(dbConn)
	// Create new server instance
	server, err := server.NewServer(store, config)
	if err != nil {
		log.Fatal("could't spin up server: %w", err)
	}
	// Setup router
	server.SetupRouter()
	// Start Server
	if err := server.StartServer("0.0.0.0:8080"); err != nil {
		log.Fatal("couldn't start server", err)
	}
}
