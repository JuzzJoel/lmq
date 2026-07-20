package middleware

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// AdminAuthMiddleware securely restricts access using a hashed secret token comparison.
func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimSpace(r.Header.Get("X-Admin-Token"))
		token = strings.ToLower(token)
		
		// Read and sanitize the environment variable hash
		envHash := strings.TrimSpace(os.Getenv("ADMIN_SECRET_HASH"))
		envHash = strings.Trim(envHash, "\"") // Strip literal double quotes
		envHash = strings.Trim(envHash, "'")  // Strip literal single quotes
		envHash = strings.ToLower(envHash)

		// Add a temporary console print to debug locally
		fmt.Printf("[AUTH DIAGNOSTIC] Expected: %s | Got: %s\n", envHash, token)

		if len(token) != 64 || len(envHash) != 64 || subtle.ConstantTimeCompare([]byte(token), []byte(envHash)) != 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: Invalid Admin Token System Key"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
