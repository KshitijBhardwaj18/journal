package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/KshitijBhardwaj18/journal/server/db"
	"github.com/KshitijBhardwaj18/journal/server/internal/auth"
	"github.com/KshitijBhardwaj18/journal/server/internal/journal"
	"github.com/KshitijBhardwaj18/journal/server/internal/user"
)

func main() {
	conn := db.NewDB()
	defer conn.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Initialize handlers
	jh := journal.NewHandler(conn)
	uh := user.NewHandler(conn)

	// Public routes
	r.Post("/auth/register", uh.Register)
	r.Post("/auth/login", uh.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		
		r.Get("/journals", jh.List)
		r.Post("/journals", jh.Create)
		r.Get("/journals/{id}", jh.GetByID)
		r.Put("/journals/{id}", jh.Update)
		r.Delete("/journals/{id}", jh.Delete)
	})

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}