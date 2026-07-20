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

// AdminAuthMiddleware securely restricts access using a hashed secret token comparison.
func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimSpace(r.Header.Get("X-Admin-Token"))
		
		// Hash the incoming token using SHA-256
		hasher := sha256.New()
		hasher.Write([]byte(token))
		hashedToken := fmt.Sprintf("%x", hasher.Sum(nil))

		// Read and sanitize the environment variable ADMIN_SECRET
		secret := strings.TrimSpace(os.Getenv("ADMIN_SECRET"))
		secret = strings.Trim(secret, "\"") // Strip literal double quotes
		secret = strings.Trim(secret, "'")  // Strip literal single quotes
		
		// Hash the expected secret
		secretHasher := sha256.New()
		secretHasher.Write([]byte(secret))
		hashedSecret := fmt.Sprintf("%x", secretHasher.Sum(nil))

		// Add a temporary console print to debug locally
		fmt.Printf("[AUTH DIAGNOSTIC] Expected: %s | Got: %s\n", hashedSecret, hashedToken)

		if subtle.ConstantTimeCompare([]byte(hashedToken), []byte(hashedSecret)) != 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: Invalid Admin Token System Key"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

