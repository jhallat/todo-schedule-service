package database

import (
	"database/sql"
	"log"
)

var DbConnection *sql.DB

func SetupDatabase(connection string) {
	var err error
	DbConnection, err = sql.Open("pgx", connection)
	if err != nil {
		log.Fatal(err)
	}
}