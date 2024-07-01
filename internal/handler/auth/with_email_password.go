package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"backendgo/internal/db"
	"backendgo/internal/response"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleLoginWithEmailPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	var user struct {
		ID       uint64 `db:"id"`
		Password string `db:"password"`
	}
	err := db.DB.QueryRow("SELECT id, password FROM identities WHERE email = $1", req.Email).Scan(&user.ID, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusUnauthorized)
			response.WriteJSON(w, map[string]string{
				"message": "invalid login credentials",
			}, nil)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response.WriteJSON(w, map[string]string{
			"message": "internal error",
		}, nil)
		slog.Error("query user info failed", slog.Any("error", err))
		return
	}

	valid := ValidatePassword(req.Password, user.Password)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		response.WriteJSON(w, map[string]string{
			"message": "invalid login credentials",
		}, nil)
		slog.Error("password is invalid", slog.String("email", req.Email))
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: strconv.FormatUint(user.ID, 10),
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(time.Hour * 24),
		},
	})

	token, err := claims.SignedString([]byte("key"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.WriteJSON(w, map[string]string{
			"message": "internal error",
		}, nil)
		slog.Error("sign token failed", slog.Any("error", err), slog.String("email", req.Email))
		return
	}

	w.WriteHeader(http.StatusOK)
	response.WriteJSON(w, map[string]any{
		"token": token,
	}, nil)
}

func ValidatePassword(password, hash string) bool {
	bytePassword := []byte(password)
	bytePasswordHash := []byte(hash)

	// comparing the password with the hash
	err := bcrypt.CompareHashAndPassword(bytePasswordHash, bytePassword)

	// nil means it is a match
	return err == nil
}
