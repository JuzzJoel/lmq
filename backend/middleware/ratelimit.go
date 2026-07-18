package middleware

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/matoous/go-nanoid/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
	"sync"
)

// RateLimiter implements a sliding window rate limiter using Redis sorted sets.
func RateLimiter(rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rdb == nil {
				next.ServeHTTP(w, r)
				return
			}
			ip := extractClientIP(r)
			key := fmt.Sprintf("ratelimit:%s", ip)
			now := time.Now()
			windowStart := now.Add(-60 * time.Second).UnixNano()
			nowNano := now.UnixNano()

			pipe := rdb.Pipeline()
			pipe.ZRemRangeByScore(r.Context(), key, "0", fmt.Sprintf("%d", windowStart))
			cardCmd := pipe.ZCard(r.Context(), key)

			if _, err := pipe.Exec(r.Context()); err != nil {
				fmt.Printf("[RateLimiter]: Redis error: %v\n", err)
				next.ServeHTTP(w, r)
				return
			}

			count := cardCmd.Val()
			if count >= 10 {
				w.Header().Set("Retry-After", "60")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Too many requests. Please try again later.",
				})
				return
			}

			randID, _ := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz", 6)
			member := fmt.Sprintf("%d-%s", nowNano, randID)

			pipe2 := rdb.Pipeline()
			pipe2.ZAdd(r.Context(), key, redis.Z{
				Score:  float64(nowNano),
				Member: member,
			})
			pipe2.Expire(r.Context(), key, 60*time.Second)
			_, _ = pipe2.Exec(r.Context())

			next.ServeHTTP(w, r)
		})
	}
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

func init() {
	go cleanupVisitors()
}

func cleanupVisitors() {
	for {
		time.Sleep(3 * time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// LocalRateLimiter implements a pure IP-based local rate limiter using golang.org/x/time/rate.
// Configured to 20 requests per minute (rate = 20/60 = 0.33, burst = 20).
func LocalRateLimiter() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractClientIP(r)

			mu.Lock()
			v, exists := visitors[ip]
			if !exists {
				limiter := rate.NewLimiter(rate.Limit(20.0/60.0), 20)
				visitors[ip] = &visitor{limiter, time.Now()}
				v = visitors[ip]
			}
			v.lastSeen = time.Now()
			mu.Unlock()

			if !v.limiter.Allow() {
				w.Header().Set("Retry-After", "60")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Too many requests. Please try again later.",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
