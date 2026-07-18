package services

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mileusna/useragent"
)

// AnalyticsService handles recording click events for links.
type AnalyticsService struct {
	pool  *pgxpool.Pool
	geoip *GeoIPService
}

// NewAnalyticsService creates a new AnalyticsService.
func NewAnalyticsService(pool *pgxpool.Pool, geoip *GeoIPService) *AnalyticsService {
	return &AnalyticsService{
		pool:  pool,
		geoip: geoip,
	}
}

// RecordClick asynchronously records a click event.
func (s *AnalyticsService) RecordClick(ctx context.Context, linkID int64, r *http.Request) {
	bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ip := extractIP(r)
	city, region, countryCode := s.geoip.LookupLocation(ip)
	uaStr := r.UserAgent()
	ua := useragent.Parse(uaStr)
	referer := r.Referer()

	tx, err := s.pool.Begin(bgCtx)
	if err != nil {
		log.Printf("[AnalyticsService]: failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback(bgCtx)

	_, err = tx.Exec(bgCtx, `
		INSERT INTO click_events (link_id, ip_address, city, region, country_code, user_agent, browser, os, is_mobile, referer)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, linkID, ip, city, region, countryCode, uaStr, ua.Name, ua.OS, ua.Mobile, referer)
	if err != nil {
		log.Printf("[AnalyticsService]: failed to insert click event: %v", err)
		return
	}

	_, err = tx.Exec(bgCtx, `
		UPDATE links SET click_count = click_count + 1 WHERE id = $1
	`, linkID)
	if err != nil {
		log.Printf("[AnalyticsService]: failed to update click count: %v", err)
		return
	}

	if err := tx.Commit(bgCtx); err != nil {
		log.Printf("[AnalyticsService]: failed to commit analytics transaction: %v", err)
	}
}

func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
