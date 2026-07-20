package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
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

// routeSpecShort is a minimal inline type for resolving A/B targets in the verify path.
type routeSpecShort struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}

func resolveTargetURLSimple(longURL string, routes []routeSpecShort) string {
	if len(routes) > 0 {
		totalWeight := 0
		for _, r := range routes {
			totalWeight += r.Weight
		}
		if totalWeight > 0 {
			roll := rand.Intn(totalWeight)
			cumulative := 0
			for _, r := range routes {
				cumulative += r.Weight
				if roll < cumulative {
					return r.URL
				}
			}
		}
	}
	return longURL
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Check if the link is burn-after-reading (needs atomic consumption)
	var isBAR bool
	var longURL string
	var hash *string
	var routesRaw []byte
	err := h.pool.QueryRow(ctx,
		"SELECT long_url, password_hash, routes, is_burn_after_reading FROM links WHERE token = $1", req.Token,
	).Scan(&longURL, &hash, &routesRaw, &isBAR)

	if err != nil || hash == nil {
		writeError(w, http.StatusNotFound, "Link not found or no password required")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(*hash), []byte(req.Password))
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	if isBAR {
		// Atomically consume the BAR link
		var deletedLongURL string
		var deletedRoutesRaw []byte
		err := h.pool.QueryRow(ctx,
			`DELETE FROM links WHERE token = $1 AND is_burn_after_reading = true
			 RETURNING long_url, routes`, req.Token,
		).Scan(&deletedLongURL, &deletedRoutesRaw)
		if err != nil {
			writeError(w, http.StatusGone, "This link has already been consumed")
			return
		}
		longURL = deletedLongURL
		routesRaw = deletedRoutesRaw
	}

	// Resolve A/B routes if present
	var routes []routeSpecShort
	if len(routesRaw) > 0 {
		json.Unmarshal(routesRaw, &routes)
	}
	resolvedURL := resolveTargetURLSimple(longURL, routes)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"long_url": resolvedURL,
	})
}
