package middleware

import (
	"net/http"
	"os"
	"strings"
)

// StrictCORSMiddleware enforces strict Cross-Origin Resource Sharing headers based on ALLOWED_ORIGINS.
func StrictCORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
			
			// If the platform handles a local direct request (no Origin header), let it resolve
			if origin != "" {
				isAllowed := false
				for _, allowed := range strings.Split(allowedOrigins, ",") {
					if origin == strings.TrimSpace(allowed) {
						isAllowed = true
						break
					}
				}
				if !isAllowed {
					http.Error(w, "CORS Policy Violation: Origin Unauthorized", http.StatusForbidden)
					return
				}
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Admin-Token")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			}
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
