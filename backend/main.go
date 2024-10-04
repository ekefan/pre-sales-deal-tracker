package main

import (
	"context"
	// FIXME: you're using two logging packages. Stick to the log/slog if you can.
	"log"
	"log/slog"

	"github.com/ekefan/pre-sales-deal-tracker/backend/api"
	db "github.com/ekefan/pre-sales-deal-tracker/backend/db/sqlc"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FIXME: I suggest you to use the "internal" pkg and expose only the needed things. By default, things should not be exposed. Then, you expose things whenever you need them to be public.

func main() {
	config, err := api.ReadConfigFiles(".")
	if err != nil {
		log.Fatal("unable to read environment variables ", err)
	}

	dbpool, err := pgxpool.New(context.Background(), config.DatabaseSource)
	if err != nil {
		log.Fatal("unable to create db connection pool ", err)
	}
	defer dbpool.Close()

	store := db.NewStore(dbpool)

	m, err := migrate.New(
		config.MigrationSource, // Path to your migrations
		config.DatabaseSource,
	)
	if err != nil {
		log.Fatal("unable to create migrate instance: ", err)
	}

	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("error applying migrations: ", err)
	} else if err == migrate.ErrNoChange {
		slog.Info("No new migrations to apply.")
	}

	// Start server
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	slog.Info("Starting HTTP server")
	server.StartServer()
}
