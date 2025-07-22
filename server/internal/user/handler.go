package user

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db *sql.DB
}

type User struct {
	ID  int `json:"id"`
	Username string `json:"username"`
	Email   string `json:"email"`
	Password string `json:"password,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User `json:"user"`
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db:db}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request){
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w,"Invalid request Body", http.StatusBadRequest)
		return
	}

	hashedPassword , err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	err = h.db.QueryRow("INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id", user.Username, user.Email, string(hashedPassword)).Scan(&user.ID)

	if err != nil {
		http.Error(w, "Username or email already exists", http.StatusConflict)
		return
	}

	token, err := generateJWT(user.ID, user.Username)

	if err != nil {
		http.Error(w,"Error generating token", http.StatusInternalServerError)
	}

	user.Password = ""
	response := AuthResponse{Token: token, User: user}
	json.NewEncoder(w).Encode(response)

}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user User
	var passwordHash string
	err := h.db.QueryRow(
		"SELECT id, username, email, password_hash FROM users WHERE username = $1",
		loginReq.Username,
	).Scan(&user.ID, &user.Username, &user.Email, &passwordHash)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(loginReq.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{Token: token, User: user}
	json.NewEncoder(w).Encode(response)
}

func generateJWT(userID int, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key")) 
}