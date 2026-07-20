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
				Description: "Create one or multiple shortened links (JSON). Payload: { url: string, routes?: [{ url: string, weight: int }], expires_in?: int, password?: string, custom_token?: string, burn_after_reading?: bool, tags?: string[] }",
			},
			{
				Method:      "POST",
				Path:        "/api/v1/shorten/csv",
				Description: "Bulk create shortened links from CSV file upload (multipart/form-data, field 'file'). Columns: url, custom_token, password, expires_in, burn_after_reading, tags.",
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
				Description: "Retrieve a paginated list of all created links. Supports ?tag= filter. Requires X-Admin-Token header.",
			},
			{
				Method:      "GET",
				Path:        "/api/v1/analytics/export",
				Description: "Export click analytics as CSV. Supports ?token=, ?from=, ?to= filters. Requires X-Admin-Token header.",
			},
			{
				Method:      "GET",
				Path:        "/{token}",
				Description: "Redirect to the resolved target URL. If the link has A/B routes, a weighted random route is selected; otherwise redirects to the base long_url.",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(docs)
}
