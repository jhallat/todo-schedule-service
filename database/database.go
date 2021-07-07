package database

import (
	"database/sql"
	"log"
)

var DbConnection *sql.DB

func SetupDatabase() {
	var err error
	//DbConnection, err = sql.Open("pgx", "postgres:Pass2021!@tcp(127.0.0.1:5432)/schedule")
	DbConnection, err = sql.Open("pgx", "user=postgres password=Pass2021! host=localhost port=5432 database=schedule sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}