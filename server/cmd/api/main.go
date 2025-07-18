package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yourusername/journal-app/db"
	"github.com/yourusername/journal-app/internal/journal"
)

func main() {
	conn := db.NewDB()
	defer conn.Close()

	r := chi.NewRouter()
	jh := journal.NewHandler(conn)

	r.Get("/journals", jh.List)
	r.Post("/journals", jh.Create)
	r.Get("/journals/{id}", jh.GetByID)
	r.Put("/journals/{id}", jh.Update)
	r.Delete("/journals/{id}", jh.Delete)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
