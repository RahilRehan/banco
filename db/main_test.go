package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	driverName = "postgres"
	dataSource = "postgres://bancoadmin:supersecret@localhost:5432/bancodb?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(driverName, dataSource)
	if err != nil {
		log.Fatalf("Cannot connect to db %v", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
