package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/ekefan/pre-sales-deal-tracker/backend/api"
	db "github.com/ekefan/pre-sales-deal-tracker/backend/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config, err := api.ReadConfigFiles(".")
	if err != nil {
		log.Fatal("unable to read environment variables ", err)
	}

	dbpool, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal("unable to create db connection pool ", err)
	}
	defer dbpool.Close()

	store := db.NewStore(dbpool)

	// Initialize migrate
	driver, err := postgres.WithInstance(dbpool, &postgres.Config{})
	if err != nil {
		log.Fatal("unable to create migrate driver: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations", // Path to your migrations
		"postgres",             // Database name
		driver,
	)
	if err != nil {
		log.Fatal("unable to create migrate instance: ", err)
	}

	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("error applying migrations: ", err)
	} else if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply.")
	}

	// Start server
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	slog.Info("starting HTTP server")
	server.StartServer()
}
