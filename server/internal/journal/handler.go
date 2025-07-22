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
	UserID  int    `json:"user_id"`
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db} 
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	
	rows, err := h.db.Query("SELECT id, title, content, user_id FROM journals WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var journals []Journal
	for rows.Next() {
		var j Journal
		if err := rows.Scan(&j.ID, &j.Title, &j.Content, &j.UserID); err != nil {
			continue
		}
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
	userID := r.Context().Value("user_id").(int)
	
	var j Journal
	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.db.QueryRow(
		"INSERT INTO journals (title, content, user_id) VALUES ($1, $2, $3) RETURNING id", 
		j.Title, j.Content, userID,
	).Scan(&j.ID)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	j.UserID = userID
	json.NewEncoder(w).Encode(j)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	id := chi.URLParam(r, "id")
	var j Journal
	json.NewDecoder(r.Body).Decode(&j)

	result, err := h.db.Exec("UPDATE journals SET title=$1, content=$2, user_id=$3 WHERE id=$4", j.Title, j.Content,userID, id,)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected();


	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if(rowsAffected == 0){
		http.Error(w,"Journal do not exist or nothing was updated", http.StatusAlreadyReported)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	id := chi.URLParam(r, "id")
	
	result, err := h.db.Exec("DELETE FROM journals WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if rowsAffected == 0 {
		http.Error(w, "Journal not found or access denied", http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}
