package handlers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"github.com/user/lmq/backend/models"
	"github.com/user/lmq/backend/services"
)

// CsvShortenHandler handles CSV bulk shortening requests.
type CsvShortenHandler struct {
	pool *pgxpool.Pool
	rdb  *redis.Client
}

// NewCsvShortenHandler creates a new CsvShortenHandler.
func NewCsvShortenHandler(pool *pgxpool.Pool, rdb *redis.Client) *CsvShortenHandler {
	return &CsvShortenHandler{pool: pool, rdb: rdb}
}

// HandleCsvShorten processes POST /api/v1/shorten/csv
func (h *CsvShortenHandler) HandleCsvShorten(w http.ResponseWriter, r *http.Request) {
	const maxUploadSize = 10 << 20 // 10 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeError(w, http.StatusBadRequest, "Failed to parse upload: "+err.Error())
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Missing 'file' field in multipart form")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to read CSV header row")
		return
	}

	colMap := buildColumnMap(header)
	if _, ok := colMap["url"]; !ok {
		writeError(w, http.StatusBadRequest, "CSV must have a 'url' column")
		return
	}

	var results []map[string]interface{}
	lineNum := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		lineNum++
		if err != nil {
			continue
		}

		result, err := h.processRow(r.Context(), record, colMap)
		if err != nil {
			results = append(results, map[string]interface{}{
				"error":  err.Error(),
				"line":   lineNum,
				"record": record,
			})
			continue
		}
		results = append(results, result)
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"results": results,
	})
}

type columnMap map[string]int

func buildColumnMap(header []string) columnMap {
	m := make(columnMap)
	for i, col := range header {
		m[strings.TrimSpace(strings.ToLower(col))] = i
	}
	return m
}

func (h *CsvShortenHandler) processRow(ctx context.Context, record []string, cm columnMap) (map[string]interface{}, error) {
	get := func(name string) string {
		if idx, ok := cm[name]; ok && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	rawURL := get("url")
	if rawURL == "" {
		return nil, fmt.Errorf("url is required")
	}

	if len(rawURL) > 2048 {
		return nil, fmt.Errorf("url exceeds 2048 characters")
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		rawURL = "https://" + rawURL
	}

	customToken := get("custom_token")
	token := customToken

	var expiresAt *time.Time
	if expStr := get("expires_in"); expStr != "" {
		var hours int
		if _, err := fmt.Sscanf(expStr, "%d", &hours); err == nil && hours > 0 {
			t := time.Now().Add(time.Duration(hours) * time.Hour)
			expiresAt = &t
		}
	}

	var passwordHash *string
	if pw := get("password"); pw != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		if err == nil {
			s := string(hash)
			passwordHash = &s
		}
	}

	isBAR := strings.ToLower(get("burn_after_reading")) == "true"

	tagStr := get("tags")
	var tags []string
	if tagStr != "" {
		for _, t := range strings.Split(tagStr, ",") {
			if trimmed := strings.TrimSpace(t); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	var routes []models.RouteSpec
	routesStr := get("routes")
	if routesStr != "" {
		json.Unmarshal([]byte(routesStr), &routes)
	}

	if h.pool != nil {
		if customToken != "" {
			var exists bool
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := h.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM links WHERE token = $1)", token).Scan(&exists)
			cancel()
			if err != nil {
				return nil, fmt.Errorf("database error checking token")
			}
			if exists {
				return nil, fmt.Errorf("custom token already exists: %s", token)
			}
		} else {
			for i := 0; i < 5; i++ {
				genToken, err := services.GenerateToken(6 + i)
				if err != nil {
					return nil, fmt.Errorf("failed to generate token")
				}
				var exists bool
				ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
				err = h.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM links WHERE token = $1)", genToken).Scan(&exists)
				cancel()
				if err == nil && !exists {
					token = genToken
					break
				}
				if i == 4 {
					return nil, fmt.Errorf("failed to generate unique token")
				}
			}
		}

		var routesJSON interface{}
		if len(routes) > 0 {
			rj, err := json.Marshal(routes)
			if err != nil {
				return nil, fmt.Errorf("failed to encode routes")
			}
			routesJSON = string(rj)
		}

		var link models.Link
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err = h.pool.QueryRow(ctx, `
			INSERT INTO links (token, long_url, expires_at, password_hash, routes, is_burn_after_reading, tags)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, token, long_url, created_at, expires_at, click_count, (password_hash IS NOT NULL), routes, is_burn_after_reading, tags
		`, token, rawURL, expiresAt, passwordHash, routesJSON, isBAR, tags).Scan(
			&link.ID, &link.Token, &link.LongURL, &link.CreatedAt, &link.ExpiresAt,
			&link.ClickCount, &link.HasPassword, &link.Routes, &link.BurnAfterReading, &link.Tags,
		)
		cancel()

		if err != nil {
			return nil, fmt.Errorf("failed to insert link: %v", err)
		}

		cacheKey := fmt.Sprintf("url:%s", token)
		h.rdb.Set(context.Background(), cacheKey, link.LongURL, 24*time.Hour)
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://lmq.name.ng"
	}

	item := map[string]interface{}{
		"token":              token,
		"short_url":          fmt.Sprintf("%s/%s", baseURL, token),
		"long_url":           rawURL,
		"created_at":         time.Now(),
		"has_password":       passwordHash != nil,
		"burn_after_reading": isBAR,
	}
	if len(tags) > 0 {
		item["tags"] = tags
	}
	if len(routes) > 0 {
		item["routes"] = routes
	}
	return item, nil
}
