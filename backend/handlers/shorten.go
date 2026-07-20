package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"github.com/user/lmq/backend/models"
	"github.com/user/lmq/backend/services"
)

var customTokenRegex = regexp.MustCompile(`^[a-zA-Z0-9-]{3,20}$`)

// ShortenHandler handles URL shortening requests.
type ShortenHandler struct {
	pool *pgxpool.Pool
	rdb  *redis.Client
}

// NewShortenHandler creates a new ShortenHandler.
func NewShortenHandler(pool *pgxpool.Pool, rdb *redis.Client) *ShortenHandler {
	return &ShortenHandler{
		pool: pool,
		rdb:  rdb,
	}
}

// HandleShorten processes POST /api/v1/shorten
func (h *ShortenHandler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.URL == "" {
		writeError(w, http.StatusBadRequest, "URL is required")
		return
	}

	rawURLs := strings.FieldsFunc(req.URL, func(c rune) bool {
		return c == ' ' || c == ','
	})

	var validURLs []string
	for _, u := range rawURLs {
		trimmed := strings.TrimSpace(u)
		if trimmed != "" {
			validURLs = append(validURLs, trimmed)
		}
	}

	if len(validURLs) == 0 {
		writeError(w, http.StatusBadRequest, "No valid URLs provided")
		return
	}

	if len(validURLs) > 50 {
		writeError(w, http.StatusBadRequest, "Bulk shortening capped at a maximum of 50 URLs per request.")
		return
	}

	if req.CustomToken != "" && len(validURLs) > 1 {
		writeError(w, http.StatusBadRequest, "Custom aliases cannot be assigned to bulk URL requests.")
		return
	}

	if req.Routes != nil && len(req.Routes) > 0 && len(validURLs) > 1 {
		writeError(w, http.StatusBadRequest, "A/B routes cannot be assigned to bulk URL requests.")
		return
	}

	if req.Routes != nil && len(req.Routes) > 0 {
		totalWeight := 0
		seen := make(map[string]bool)
		for _, r := range req.Routes {
			if r.URL == "" {
				writeError(w, http.StatusBadRequest, "Each route must have a URL.")
				return
			}
			if r.Weight <= 0 {
				writeError(w, http.StatusBadRequest, "Each route must have a positive weight.")
				return
			}
			if seen[r.URL] {
				writeError(w, http.StatusBadRequest, "Duplicate route URLs are not allowed.")
				return
			}
			seen[r.URL] = true
			totalWeight += r.Weight
		}
		if totalWeight <= 0 {
			writeError(w, http.StatusBadRequest, "Route weights must sum to a positive value.")
			return
		}
	}

	var results []map[string]interface{}

	if h.pool == nil {
		for i, vurl := range validURLs {
			parsedURL, err := url.ParseRequestURI(vurl)
			if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
				vurl = "https://" + vurl
			}
			mockToken := fmt.Sprintf("mockB62-%d", i)
			baseURL := os.Getenv("BASE_URL")
			if baseURL == "" {
				baseURL = "https://lmq.name.ng"
			}
			item := map[string]interface{}{
				"short_url":          fmt.Sprintf("%s/%s", baseURL, mockToken),
				"token":              mockToken,
				"long_url":           vurl,
				"created_at":         time.Now(),
				"burn_after_reading": req.BurnAfterReading,
			}
			if len(req.Routes) > 0 {
				item["routes"] = req.Routes
			}
			results = append(results, item)
		}
		writeJSON(w, http.StatusCreated, map[string]interface{}{
			"results": results,
			"mock":    true,
		})
		return
	}

	for _, vurl := range validURLs {
		if len(vurl) > 2048 {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("URL exceeds maximum length of 2048 characters: %s", vurl))
			return
		}
		parsedURL, err := url.ParseRequestURI(vurl)
		if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
			vurl = "https://" + vurl
		}

		var token string
		if req.CustomToken != "" {
			token = req.CustomToken
			if !customTokenRegex.MatchString(token) {
				writeError(w, http.StatusBadRequest, "Custom token must be 3-20 characters long and contain only alphanumeric characters and hyphens")
				return
			}
			var exists bool
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			err := h.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM links WHERE token = $1)", token).Scan(&exists)
			cancel()
			if err != nil {
				writeError(w, http.StatusInternalServerError, "Database error")
				return
			}
			if exists {
				writeError(w, http.StatusConflict, "Custom token already exists")
				return
			}
		} else {
			for i := 0; i < 5; i++ {
				genToken, err := services.GenerateToken(6 + i)
				if err != nil {
					writeError(w, http.StatusInternalServerError, "Failed to generate token")
					return
				}
				var exists bool
				ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
				err = h.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM links WHERE token = $1)", genToken).Scan(&exists)
				cancel()
			if err == nil && !exists {
					token = genToken
					break
				}
				if i == 4 {
					writeError(w, http.StatusInternalServerError, "Failed to generate a unique token after 5 attempts")
					return
				}
			}
		}

		var expiresAt *time.Time
		if req.ExpiresIn > 0 {
			t := time.Now().Add(time.Duration(req.ExpiresIn) * time.Hour)
			expiresAt = &t
		}

		var passwordHash *string
		if req.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "Failed to secure password")
				return
			}
			hStr := string(hash)
			passwordHash = &hStr
		}

		var link models.Link
		var routesJSON interface{}
		if len(req.Routes) > 0 {
			rj, err := json.Marshal(req.Routes)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "Failed to encode routes")
				return
			}
			routesJSON = string(rj)
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		err = h.pool.QueryRow(ctx, `
			INSERT INTO links (token, long_url, expires_at, password_hash, routes, is_burn_after_reading) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING id, token, long_url, created_at, expires_at, click_count, (password_hash IS NOT NULL), routes, is_burn_after_reading
		`, token, vurl, expiresAt, passwordHash, routesJSON, req.BurnAfterReading).Scan(&link.ID, &link.Token, &link.LongURL, &link.CreatedAt, &link.ExpiresAt, &link.ClickCount, &link.HasPassword, &link.Routes, &link.BurnAfterReading)
		cancel()

		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to insert link into database")
			return
		}

		cacheKey := fmt.Sprintf("url:%s", token)
		h.rdb.Set(context.Background(), cacheKey, link.LongURL, 24*time.Hour)

		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "https://lmq.name.ng"
		}
		shortURL := fmt.Sprintf("%s/%s", baseURL, link.Token)

		resultItem := map[string]interface{}{
			"token":              link.Token,
			"short_url":          shortURL,
			"long_url":           link.LongURL,
			"created_at":         link.CreatedAt,
			"expires_at":         link.ExpiresAt,
			"has_password":       link.HasPassword,
			"burn_after_reading": link.BurnAfterReading,
		}
		if len(link.Routes) > 0 {
			resultItem["routes"] = link.Routes
		}
		results = append(results, resultItem)
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"results": results,
	})
}
