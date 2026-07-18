package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type VerifyHandler struct {
	pool *pgxpool.Pool
}

func NewVerifyHandler(pool *pgxpool.Pool) *VerifyHandler {
	return &VerifyHandler{pool: pool}
}

type VerifyRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (h *VerifyHandler) HandleVerify(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.Token == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "Token and password required")
		return
	}

	if h.pool == nil {
		writeError(w, http.StatusInternalServerError, "Database unavailable")
		return
	}

	var longURL string
	var hash *string
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := h.pool.QueryRow(ctx, "SELECT long_url, password_hash FROM links WHERE token = $1", req.Token).Scan(&longURL, &hash)
	if err != nil || hash == nil {
		writeError(w, http.StatusNotFound, "Link not found or no password required")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(*hash), []byte(req.Password))
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"long_url": longURL,
	})
}
