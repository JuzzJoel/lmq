package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/user/lmq/backend/models"
)

// AnalyticsHandler serves analytic data for short links.
type AnalyticsHandler struct {
	pool *pgxpool.Pool
}

// NewAnalyticsHandler constructs a handler for analytics API endpoints.
func NewAnalyticsHandler(pool *pgxpool.Pool) *AnalyticsHandler {
	return &AnalyticsHandler{pool: pool}
}

// HandleGetAnalytics returns detailed analytics for a single link token.
// Endpoint: GET /api/v1/analytics?token=xxxxx
func (h *AnalyticsHandler) HandleGetAnalytics(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		writeError(w, http.StatusBadRequest, "Query parameter 'token' is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if h.pool == nil {
		writeJSON(w, http.StatusOK, models.LinkAnalytics{
			Token:         token,
			LongURL:       "https://mock-backend.local",
			TotalClicks:   999,
			ClicksByDay:   []models.DayCount{{Date: time.Now().Format("2006-01-02"), Count: 999}},
			Cities:        []models.CityCount{{City: "Mock City", Count: 999}},
			Regions:       []models.RegionCount{{Region: "Mock Region", Count: 999}},
			CountryGroups: []models.CountryCount{{CountryCode: "XX", Count: 999}},
			Browsers:      []models.BrowserCount{{Browser: "MockBrowser", Count: 999}},
			RecentClicks:  []models.ClickEvent{},
		})
		return
	}

	var linkID int64
	var longURL string
	var totalClicks int64
	err := h.pool.QueryRow(ctx,
		"SELECT id, long_url, click_count FROM links WHERE token = $1", token,
	).Scan(&linkID, &longURL, &totalClicks)
	if err != nil {
		writeError(w, http.StatusNotFound, "Link not found")
		return
	}

	clicksByDay := make([]models.DayCount, 0)
	rows, err := h.pool.Query(ctx, `
		SELECT DATE(clicked_at) AS day, COUNT(*) AS cnt
		FROM click_events
		WHERE link_id = $1 AND clicked_at >= NOW() - INTERVAL '30 days'
		GROUP BY day
		ORDER BY day ASC
	`, linkID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var dc models.DayCount
			var dayTime time.Time
			if scanErr := rows.Scan(&dayTime, &dc.Count); scanErr == nil {
				dc.Date = dayTime.Format("2006-01-02")
				clicksByDay = append(clicksByDay, dc)
			}
		}
	}

	cities := make([]models.CityCount, 0)
	cityRows, err := h.pool.Query(ctx, `
		SELECT city, COUNT(*) AS cnt
		FROM click_events
		WHERE link_id = $1
		GROUP BY city
		ORDER BY cnt DESC
		LIMIT 20
	`, linkID)
	if err == nil {
		defer cityRows.Close()
		for cityRows.Next() {
			var cc models.CityCount
			if scanErr := cityRows.Scan(&cc.City, &cc.Count); scanErr == nil {
				cities = append(cities, cc)
			}
		}
	}

	regions := make([]models.RegionCount, 0)
	regionRows, err := h.pool.Query(ctx, `
		SELECT region, COUNT(*) AS cnt
		FROM click_events
		WHERE link_id = $1
		GROUP BY region
		ORDER BY cnt DESC
		LIMIT 20
	`, linkID)
	if err == nil {
		defer regionRows.Close()
		for regionRows.Next() {
			var rc models.RegionCount
			if scanErr := regionRows.Scan(&rc.Region, &rc.Count); scanErr == nil {
				regions = append(regions, rc)
			}
		}
	}

	countryGroups := make([]models.CountryCount, 0)
	countryRows, err := h.pool.Query(ctx, `
		SELECT country_code, COUNT(*) AS cnt
		FROM click_events
		WHERE link_id = $1
		GROUP BY country_code
		ORDER BY cnt DESC
		LIMIT 20
	`, linkID)
	if err == nil {
		defer countryRows.Close()
		for countryRows.Next() {
			var cc models.CountryCount
			if scanErr := countryRows.Scan(&cc.CountryCode, &cc.Count); scanErr == nil {
				countryGroups = append(countryGroups, cc)
			}
		}
	}

	browsers := make([]models.BrowserCount, 0)
	browserRows, err := h.pool.Query(ctx, `
		SELECT browser, COUNT(*) AS cnt
		FROM click_events
		WHERE link_id = $1 AND browser != ''
		GROUP BY browser
		ORDER BY cnt DESC
		LIMIT 10
	`, linkID)
	if err == nil {
		defer browserRows.Close()
		for browserRows.Next() {
			var bc models.BrowserCount
			if scanErr := browserRows.Scan(&bc.Browser, &bc.Count); scanErr == nil {
				browsers = append(browsers, bc)
			}
		}
	}

	recentClicks := make([]models.ClickEvent, 0)
	recentRows, err := h.pool.Query(ctx, `
		SELECT id, link_id, clicked_at, city, region, country_code, browser, os, is_mobile, COALESCE(referer, '')
		FROM click_events
		WHERE link_id = $1
		ORDER BY clicked_at DESC
		LIMIT 50
	`, linkID)
	if err == nil {
		defer recentRows.Close()
		for recentRows.Next() {
			var ce models.ClickEvent
			if scanErr := recentRows.Scan(
				&ce.ID, &ce.LinkID, &ce.ClickedAt,
				&ce.City, &ce.Region, &ce.CountryCode, &ce.Browser, &ce.OS, &ce.IsMobile, &ce.Referer,
			); scanErr == nil {
				recentClicks = append(recentClicks, ce)
			}
		}
	}

	analytics := models.LinkAnalytics{
		Token:         token,
		LongURL:       longURL,
		TotalClicks:   totalClicks,
		ClicksByDay:   clicksByDay,
		Cities:        cities,
		Regions:       regions,
		CountryGroups: countryGroups,
		Browsers:      browsers,
		RecentClicks:  recentClicks,
	}

	writeJSON(w, http.StatusOK, analytics)
}

// HandleListLinks returns paginated shortened links with summary click counts.
// Endpoint: GET /api/v1/analytics/links?page=1&limit=20&search=keyword
func (h *AnalyticsHandler) HandleListLinks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	page := 1
	limit := 20
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}
	search := r.URL.Query().Get("search")
	offset := (page - 1) * limit

	if h.pool == nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"links": []models.Link{
				{ID: 1, Token: "mockB62-0", LongURL: "https://mock.link.1", CreatedAt: time.Now(), ClickCount: 150},
				{ID: 2, Token: "mockB62-1", LongURL: "https://mock.link.2", CreatedAt: time.Now(), ClickCount: 300},
			},
			"total": 2,
		})
		return
	}

	var totalCount int64
	var err error

	if search != "" {
		err = h.pool.QueryRow(ctx, `SELECT COUNT(*) FROM links WHERE long_url ILIKE '%' || $1 || '%' OR token ILIKE '%' || $1 || '%'`, search).Scan(&totalCount)
	} else {
		err = h.pool.QueryRow(ctx, `SELECT COUNT(*) FROM links`).Scan(&totalCount)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to count links")
		return
	}

	var rows pgx.Rows
	// Note: We need pgx import for pgx.Rows but we can just use the returned rows without defining var rows pgx.Rows.
	// Let's do it directly.
	if search != "" {
		rows, err = h.pool.Query(ctx, `
			SELECT id, token, long_url, created_at, expires_at, click_count, (password_hash IS NOT NULL)
			FROM links
			WHERE long_url ILIKE '%' || $1 || '%' OR token ILIKE '%' || $1 || '%'
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`, search, limit, offset)
	} else {
		rows, err = h.pool.Query(ctx, `
			SELECT id, token, long_url, created_at, expires_at, click_count, (password_hash IS NOT NULL)
			FROM links
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`, limit, offset)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve links")
		return
	}
	defer rows.Close()

	links := make([]models.Link, 0)
	for rows.Next() {
		var link models.Link
		if scanErr := rows.Scan(
			&link.ID, &link.Token, &link.LongURL,
			&link.CreatedAt, &link.ExpiresAt, &link.ClickCount, &link.HasPassword,
		); scanErr == nil {
			links = append(links, link)
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"links": links,
		"total": totalCount,
	})
}
