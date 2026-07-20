package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/html"

	"github.com/user/lmq/backend/models"
	"github.com/user/lmq/backend/services"
)

// RedirectHandler manages the GET /{token} redirection pipeline.
type RedirectHandler struct {
	pool      *pgxpool.Pool
	rdb       *redis.Client
	analytics *services.AnalyticsService
}

// NewRedirectHandler constructs a handler for short URL redirection.
func NewRedirectHandler(pool *pgxpool.Pool, rdb *redis.Client, analytics *services.AnalyticsService) *RedirectHandler {
	return &RedirectHandler{
		pool:      pool,
		rdb:       rdb,
		analytics: analytics,
	}
}

type cachedLink struct {
	LongURL     string           `json:"long_url"`
	ExpiresAt   *time.Time       `json:"expires_at"`
	HasPassword bool             `json:"has_password"`
	ID          int64            `json:"id"`
	Routes      []models.RouteSpec `json:"routes,omitempty"`
}

// resolveTargetURL picks a destination from routes via weighted random, or falls back to long_url.
func resolveTargetURL(link *cachedLink) string {
	if len(link.Routes) > 0 {
		totalWeight := 0
		for _, r := range link.Routes {
			totalWeight += r.Weight
		}
		if totalWeight > 0 {
			roll := rand.Intn(totalWeight)
			cumulative := 0
			for _, r := range link.Routes {
				cumulative += r.Weight
				if roll < cumulative {
					return r.URL
				}
			}
		}
	}
	return link.LongURL
}

// HandleRedirect intercepts GET /{token}, performs checks, and redirects or unfurls.
func (h *RedirectHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var link cachedLink

	if h.pool == nil {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	// Step 1: Check Database directly to ensure secure state for passwords and expiration
	// (We bypass Redis cache here for absolute source-of-truth on TTL/Passwords)
	var expiresAt *time.Time
	var passwordHash *string
	var routesRaw []byte
	err := h.pool.QueryRow(ctx,
		"SELECT id, long_url, expires_at, password_hash, routes FROM links WHERE token = $1", token,
	).Scan(&link.ID, &link.LongURL, &expiresAt, &passwordHash, &routesRaw)

	if err != nil {
		http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
		return
	}
	link.ExpiresAt = expiresAt
	link.HasPassword = passwordHash != nil

	if len(routesRaw) > 0 {
		if err := json.Unmarshal(routesRaw, &link.Routes); err != nil {
			link.Routes = nil
		}
	}

	// Resolve target URL via A/B routing or direct long_url
	targetURL := resolveTargetURL(&link)

	// Expiration check
	if link.ExpiresAt != nil && link.ExpiresAt.Before(time.Now()) {
		http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
		return
	}

	// Password check
	if link.HasPassword {
		http.Redirect(w, r, fmt.Sprintf("/protected/%s", token), http.StatusTemporaryRedirect)
		return
	}

	// Crawler Detection for OG Unfurling
	ua := strings.ToLower(r.UserAgent())
	if strings.Contains(ua, "twitterbot") || strings.Contains(ua, "discordbot") || strings.Contains(ua, "slackbot") {
		h.handleCrawlerUnfurl(w, targetURL)
		return
	}

	// Fire-and-forget analytics recording
	if link.ID > 0 {
		go h.analytics.RecordClick(context.Background(), link.ID, r)
	}

	// Issue HTTP 301 Moved Permanently
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Location", targetURL)
	w.WriteHeader(http.StatusMovedPermanently)
}

func (h *RedirectHandler) handleCrawlerUnfurl(w http.ResponseWriter, targetURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		http.Redirect(w, nil, targetURL, http.StatusMovedPermanently)
		return
	}
	req.Header.Set("User-Agent", "LMQ-Bot/1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Redirect(w, nil, targetURL, http.StatusMovedPermanently)
		return
	}
	defer resp.Body.Close()

	title, ogTags := extractMetaTags(resp.Body)
	
	if title == "" {
		title = "LMQ Redirect"
	}

	htmlOut := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
    %s
    <meta http-equiv="refresh" content="0; url=%s">
</head>
<body>
    <p>Redirecting to <a href="%s">%s</a>...</p>
</body>
</html>`, title, ogTags, targetURL, targetURL, targetURL)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlOut))
}

func extractMetaTags(body io.Reader) (string, string) {
	z := html.NewTokenizer(body)
	title := ""
	var ogBuilder strings.Builder

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			t := z.Token()

			if t.Data == "title" {
				tt = z.Next()
				if tt == html.TextToken {
					title = z.Token().Data
				}
			} else if t.Data == "meta" {
				var isOg bool
				var prop, content string
				for _, a := range t.Attr {
					if a.Key == "property" && strings.HasPrefix(a.Val, "og:") {
						isOg = true
						prop = a.Val
					} else if a.Key == "content" {
						content = a.Val
					}
				}
				if isOg {
					ogBuilder.WriteString(fmt.Sprintf(`<meta property="%s" content="%s">`+"\n", prop, content))
				}
			}
		}
	}
	return title, ogBuilder.String()
}
