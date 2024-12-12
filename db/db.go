package db

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

func Connect(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}
	return db
}
