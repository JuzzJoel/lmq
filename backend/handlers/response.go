package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/user/lmq/backend/models"
)

// writeJSON writes a generic JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.APIResponse[interface{}]{
		Data:  data,
		Error: nil,
	})
}

// writeError writes a JSON error response
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.APIResponse[interface{}]{
		Data:  nil,
		Error: &message,
	})
}
