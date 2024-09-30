package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/ekefan/pre-sales-deal-tracker/backend/api"
	db "github.com/ekefan/pre-sales-deal-tracker/backend/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), "postgresql://root:vasDealTracker@localhost:5432/dealTrackerDB?sslmode=disable")
	if err != nil {
		log.Fatal("Unable to create connection pool", err)
	}
	defer dbpool.Close()

	store := db.NewStore(dbpool)
	server := api.NewServer(store)
	slog.Info("starting http server")
	server.StartServer(":8080")
}