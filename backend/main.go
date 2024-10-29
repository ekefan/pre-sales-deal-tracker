package main

import (
	"context"
	"os"
	"log/slog"

	"github.com/ekefan/pre-sales-deal-tracker/backend/api"
	db "github.com/ekefan/pre-sales-deal-tracker/backend/internal/db/sqlc"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := api.ReadConfigFiles(".")
	if err != nil {
		slog.Error("unable to read environment variables", "error", err)
		os.Exit(1)
	}
	dbpool, err := pgxpool.New(context.Background(), config.DatabaseSource)
	if err != nil {
		slog.Error("unable to create db connection pool ", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	store := db.NewStore(dbpool)
	m, err := migrate.New(
		config.MigrationSource, // Path to your migrations
		config.DatabaseSource,
	)
	if err != nil {
		slog.Error("unable to create migrate instance", "error", err)
		os.Exit(1)
	}
	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("error applying migrations: ", "error", err)
		os.Exit(1)
	} else if err == migrate.ErrNoChange {
		slog.Info("No new migrations to apply.")
	}
	// Start server
	server, err := api.NewServer(store, config)
	if err != nil {
		slog.Error("cannot start server: ", "error", err)
	}
	slog.Info("Starting HTTP server")
	server.StartServer()
}
