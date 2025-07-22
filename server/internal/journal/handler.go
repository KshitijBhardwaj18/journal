package journal

import (
	"database/sql"
	"encoding/json"
	"net/http"
	// "strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	db *sql.DB
}

type Journal struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db} 
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	rows, _ := h.db.Query("SELECT id, title, content FROM journals")
	defer rows.Close()

	var journals []Journal
	for rows.Next() {
		var j Journal
		rows.Scan(&j.ID, &j.Title, &j.Content)
		journals = append(journals, j)
	}
	json.NewEncoder(w).Encode(journals)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	row := h.db.QueryRow("SELECT id, title, content FROM journals WHERE id=$1", id)

	var j Journal
	if err := row.Scan(&j.ID, &j.Title, &j.Content); err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(j)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var j Journal
	json.NewDecoder(r.Body).Decode(&j)

	err := h.db.QueryRow("INSERT INTO journals (title, content) VALUES ($1, $2) RETURNING id", j.Title, j.Content).Scan(&j.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(j)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var j Journal
	json.NewDecoder(r.Body).Decode(&j)

	_, err := h.db.Exec("UPDATE journals SET title=$1, content=$2 WHERE id=$3", j.Title, j.Content, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec("DELETE FROM journals WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
