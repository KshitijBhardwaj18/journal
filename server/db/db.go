package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewDB() *sql.DB {
	connStr := "postgres://postgres:yourpassword@localhost:5432/journaldb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
