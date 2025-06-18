package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/dohoanggiahuy317/ACH-transactions-Microservice-app/db/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Initialize the database connection
	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Initialize the test queries
	testQueries = New(testDB)

	// Run the tests
	os.Exit(m.Run())
}