package auth

import (
	"backendgo/internal/db"
	"backendgo/internal/response"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	var found bool
	err := db.DB.QueryRow("SELECT 1 FROM identities WHERE email = $1", req.Email).Scan(&found)
	if found {
		w.WriteHeader(http.StatusBadRequest)
		response.WriteJSON(w, map[string]string{
			"message": "the given email is existing, please login or choose another email",
		}, nil)
		log.Printf("email %s existing", req.Email)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "internal error",
		})
		log.Printf("cannot generate hash %v", err)
		return
	}
	_, err = db.DB.Exec(`INSERT INTO identities (email, password) VALUES ($1, $2)`, req.Email, hash)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "register succeeded",
	})
}
