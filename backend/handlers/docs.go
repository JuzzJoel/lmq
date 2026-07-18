package handlers

import (
	"encoding/json"
	"net/http"
)

type APIEndpoint struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type APIDocsResponse struct {
	Version   string        `json:"version"`
	Endpoints []APIEndpoint `json:"endpoints"`
}

func HandleGetDocs(w http.ResponseWriter, r *http.Request) {
	docs := APIDocsResponse{
		Version: "1.0.0",
		Endpoints: []APIEndpoint{
			{
				Method:      "POST",
				Path:        "/api/v1/shorten",
				Description: "Create one or multiple shortened links. Payload: { url: string, expires_in?: int, password?: string, alias?: string }",
			},
			{
				Method:      "POST",
				Path:        "/api/v1/verify-password",
				Description: "Unlock a password-protected link before redirecting. Payload: { short_id: string, password: string }",
			},
			{
				Method:      "GET",
				Path:        "/api/v1/analytics",
				Description: "Retrieve global link performance analytics including cities, regions, and country groups. Requires X-Admin-Token header.",
			},
			{
				Method:      "GET",
				Path:        "/api/v1/analytics/links",
				Description: "Retrieve a paginated list of all created links. Requires X-Admin-Token header.",
			},
			{
				Method:      "GET",
				Path:        "/{token}",
				Description: "Redirect to the original long URL associated with this short token.",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(docs)
}
