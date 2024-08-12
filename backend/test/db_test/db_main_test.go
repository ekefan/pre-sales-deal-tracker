package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
)

const (
	dbDriver string = "postgres"
	dbSource string = "postgresql://root:vasDealTracker@localhost:5432/dealTrackerDB?sslmode=disable"
)

var ts db.Store

func TestMain(m *testing.M) {
	dbConn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("couldn't connect to database")
	}
	ts = db.NewStore(dbConn)

	os.Exit(m.Run())
}
