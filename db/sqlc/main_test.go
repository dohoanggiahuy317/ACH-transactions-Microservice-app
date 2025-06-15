package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres" // or your database driver
	dbSource = "postgresql://root:secret@localhost:55432/simple_bank?sslmode=disable" // Update with your connection string
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	// Initialize the database connection
	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Initialize the test queries
	testQueries = New(testDB)

	// Run the tests
	os.Exit(m.Run())
}