package main

import (
	"database/sql"
	"log"

	"github.com/dohoanggiahuy317/ACH-transactions-Microservice-app/api"
	db "github.com/dohoanggiahuy317/ACH-transactions-Microservice-app/db/sqlc"
	"github.com/dohoanggiahuy317/ACH-transactions-Microservice-app/db/util"
	_ "github.com/lib/pq"
)


func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Initialize the database connection
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// Verify the database connection
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
