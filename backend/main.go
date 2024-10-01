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
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("can not start server, ", err)
	}
	slog.Info("starting http server")
	server.StartServer()
}