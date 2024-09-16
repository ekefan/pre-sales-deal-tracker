package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"fmt"
	"math/rand"
	"strings"
	_ "github.com/lib/pq"
)

const (
	dbDriver string = "postgres"
	dbSource string = "postgresql://root:vasDealTracker@localhost:5432/dealTrackerDB?sslmode=disable"
)

var ts Store

func TestMain(m *testing.M) {
	dbConn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("couldn't connect to database")
	}
	ts = NewStore(dbConn)

	os.Exit(m.Run())
}




var alphabets = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		char := alphabets[rand.Intn(len(alphabets))]
		sb.WriteByte(char)
	}
	return sb.String()
}


func GenEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(5))
}

func GenFullname() string {
	return fmt.Sprintf("%s %s", RandomString(5), RandomString(4))
}

func GenPassWord() string {
	return fmt.Sprintf("%s%v", RandomString(6), 10 + rand.Int63n(89))
}
