package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <migration_file>")
	}

	connStr := "postgres://postgres:cooldude@localhost:5434/journaldb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrationFile := os.Args[1]
	content, err := ioutil.ReadFile(migrationFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Migration %s applied successfully!\n", migrationFile)
}