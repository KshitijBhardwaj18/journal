package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewDB() *sql.DB {
	connStr := "postgres://postgres:cooldude@localhost:5434/journaldb?sslmode=disable"
	log.Printf("Attempting to connect to database with: %s", connStr)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		log.Fatal(err)
	}
	
	log.Println("Database connection opened successfully, testing ping...")
	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		log.Fatal(err)
	}
	
	log.Println("Database connection successful!")
	return db
}
