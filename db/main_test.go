package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const(
	driverName = "postgres"
	dataSource = "postgres://bancoadmin:supersecret@localhost:5432/bancodb?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return
	}

	testQueries = New(db)
	os.Exit(m.Run())
}