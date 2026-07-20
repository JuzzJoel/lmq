package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// AdminAuthMiddleware securely restricts access using a pre-computed SHA-256 hash comparison.
func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimSpace(r.Header.Get("X-Admin-Token"))

		// Hash the incoming token using SHA-256
		hasher := sha256.New()
		hasher.Write([]byte(token))
		hashedToken := fmt.Sprintf("%x", hasher.Sum(nil))

		// Read the pre-computed SHA-256 hash from ADMIN_SECRET_HASH (not plaintext)
		expectedHash := strings.TrimSpace(os.Getenv("ADMIN_SECRET_HASH"))
		expectedHash = strings.Trim(expectedHash, "\"") // Strip literal double quotes
		expectedHash = strings.Trim(expectedHash, "'")  // Strip literal single quotes

		if subtle.ConstantTimeCompare([]byte(hashedToken), []byte(expectedHash)) != 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: Invalid Admin Token System Key"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

